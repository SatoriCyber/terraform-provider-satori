package resources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"os"
)

func ResourceUserCustomAttributes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCustomAttributesCreate,
		ReadContext:   resourceUserCustomAttributesRead,
		UpdateContext: resourceUserCustomAttributesUpdate,
		DeleteContext: resourceUserCustomAttributesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Satori user settings allows to config existing user's configuration. Currently supports only user's attributes configuration",
		Schema:      getUserCustomAttributesDefinitionSchema(),
	}
}

func resourceUserCustomAttributesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	_, err := createUserAttributes(d, c)

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func createUserAttributes(d *schema.ResourceData, c *api.Client) (*api.UserWithCustomAttributes, error) {
	input, err := resourceToUserAttribute(d)

	if err != nil {
		return nil, err
	}

	result, err := c.CreateUserCustomAttributes(input)

	if err != nil {
		return nil, err
	}
	d.SetId(d.Get("user_id").(string))

	return result, err
}

func resourceToUserAttribute(d *schema.ResourceData) (*api.UserWithCustomAttributes, error) {
	attrDto := api.UserWithCustomAttributes{}

	if d == nil {
		return nil, nil
	}

	attrRawJson := d.Get(Attributes).(string)
	attrDto.UserId = d.Get(UserId).(string)

	fileContent, readFileErr := os.ReadFile(attrRawJson)

	if readFileErr == nil {
		attrRawJson = string(fileContent)
	}

	err := json.Unmarshal([]byte(attrRawJson), &attrDto.CustomAttributes)

	if err != nil {
		return nil, err
	}

	if !validMapElementsAttributesType(attrDto.CustomAttributes) {
		errMsg := fmt.Sprintf("Each attribute element in the list must be one of the following list: { string, int, float, bool, []string, []number } where number is int|float. This is not the case for the resource defined for user with ID '%s'", d.Get(UserId))
		return nil, errors.New(errMsg)
	}

	return &attrDto, nil
}

func resourceUserCustomAttributesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*api.Client)

	_, err := getUserAttributes(c, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func getUserAttributes(c *api.Client, d *schema.ResourceData) (*api.UserWithCustomAttributes, error) {

	userId := (*d).Get(UserId).(string)

	result, err, responseStatus := c.GetUserProfile(&userId)

	if responseStatus == 404 || (*result).CustomAttributes == nil {
		return nil, errors.New(fmt.Sprintf("User with id '%s' does not exists.", userId))
	}

	if err != nil {
		return nil, err
	}

	mergedAttributes, err := mergeUserAndConfiguredAttributesMap((*result).CustomAttributes, d)
	if err != nil {
		return nil, err
	}

	rawMergedAttributes, err := json.Marshal(mergedAttributes)
	if err != nil {
		return nil, err
	}

	err = d.Set(Attributes, string(rawMergedAttributes))
	if err != nil {
		return nil, err
	}
	d.SetId((*result).UserId)

	if err != nil {
		return nil, err
	}

	attr := api.UserWithCustomAttributes{}

	attr.UserId = d.Get(UserId).(string)
	attr.CustomAttributes = mergedAttributes

	return &attr, err
}

// This function merges the attributes from the user resource retrieved from backend and the ResourceData schema.
func mergeUserAndConfiguredAttributesMap(userAttrMap map[string]interface{}, d *schema.ResourceData) (map[string]interface{}, error) {
	changeMap := make(map[string]interface{})
	currAttributesMap := make(map[string]interface{})
	rawAttributes := d.Get(Attributes).(string)

	fileContent, readFileErr := os.ReadFile(rawAttributes)

	if readFileErr == nil {
		rawAttributes = string(fileContent)
	}

	err := json.Unmarshal([]byte(rawAttributes), &currAttributesMap)
	if err != nil {
		return nil, err
	}

	for key := range currAttributesMap {
		// Keeping change map updated in-case there is a diff between currAttributesMap and userAttrMap
		// which means that the attributes should be updated.
		userVal, _ := userAttrMap[key]
		changeMap[key] = userVal

		// Delete the key, so it will not be duplicated down the road.
		delete(userAttrMap, key)
	}

	// If there are any attributes left in the attributes set.
	// Add them to the changeMap
	for key, userValue := range userAttrMap {
		changeMap[key] = userValue
	}

	return changeMap, nil
}

func resourceUserCustomAttributesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	_, err := updateSatoriAttributes(d, c)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateSatoriAttributes(d *schema.ResourceData, c *api.Client) (*api.UserWithCustomAttributes, error) {
	input, err := resourceToUserAttribute(d)

	if err != nil {
		return nil, err
	}
	result, err := c.UpdateUserCustomAttributes(input)
	return result, err
}

func resourceUserCustomAttributesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	_, err := deleteUserSatoriAttributes(d, c)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteUserSatoriAttributes(d *schema.ResourceData, c *api.Client) (*api.UserWithCustomAttributes, error) {
	input := api.UserWithCustomAttributes{}
	input.UserId = d.Get(UserId).(string)

	result, err := c.DeleteUserCustomAttributes(&input)

	// Setting terraform-id from the empty string ---> triggers delete
	d.SetId("")

	return result, err
}
