package resources

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

func ResourceDataStore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataStoreCreate,
		ReadContext:   resourceDataStoreRead,
		UpdateContext: resourceDataStoreUpdate,
		DeleteContext: resourceDataStoreDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Full DataStore configuration.",
		Schema:      getDataStoreDefinitionSchema(),
	}
}

func resourceDataStoreCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	result, err := createDataStore(d, c)
	if err != nil {
		return diag.FromErr(err)
	}
	fmt.Println(result.Id)
	//if err = updateDataPolicy(err, c, result.DataPolicyId, d); err != nil {
	//	return diag.FromErr(err)
	//}

	return diags
}

func resourceDataStoreRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	result, err := getDataStore(c, d)
	if result == nil && err == nil {
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	//resultCustomPolicy, err := c.GetCustomPolicy(result.DataPolicyId)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//
	//if err := d.Set("custom_policy", []map[string]interface{}{*customPolicyToResource(resultCustomPolicy)}); err != nil {
	//	return diag.FromErr(err)
	//}
	//
	//resultAccessControl, err := c.GetAccessControl(result.DataPolicyId)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//
	//if err := d.Set("access_control_settings", []map[string]interface{}{*accessControlToResource(resultAccessControl)}); err != nil {
	//	return diag.FromErr(err)
	//}
	//
	//resultSecurityPolicies, err := c.GetSecurityPolicies(result.DataPolicyId)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//
	//if err := setStringListProp(&resultSecurityPolicies.Ids, "security_policies", d); err != nil {
	//	return diag.FromErr(err)
	//}

	return diags
}

func resourceDataStoreUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	result, err := updateDataStore(d, c)
	if err != nil {
		return diag.FromErr(err)
	}
	fmt.Println(result)
	//if err = updateDataPolicy(err, c, result.DataPolicyId, d); err != nil {
	//	return diag.FromErr(err)
	//}

	return diags
}
