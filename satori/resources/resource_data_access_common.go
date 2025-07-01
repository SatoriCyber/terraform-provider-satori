package resources

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"log"
)

var (
	identityNamePath    = "identity.0.name"
	identityTypePath    = "identity.0.type"
	identityGroupIdPath = "identity.0.group_id"
)

func resourceDataAccessParent() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Parent data policy ID, the data_policy_id field of a dataset.",
	}
}

func resourceDataAccessLevel() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Access level to grant, valid values are: READ_ONLY, READ_WRITE, OWNER.",
	}
}

func resourceDataAccessRevokeIfNotUsed() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     0,
		Description: "Revoke access if rule not used in the last given days. Zero = do not revoke. Max value is 180.",
	}
}

func resourceDataAccessIdentity(isCelSupported bool) *schema.Schema {
	var identityTypeDescription = fmt.Sprintf("Identity type, valid types are: %s", validIdentityTypeList(isCelSupported))

	var identityNameDescription = "User/group name for identity types of USER and IDP_GROUP"
	if isCelSupported {
		identityNameDescription += " or a custom expression based on attributes of the identity for CEL identity type"
	}
	identityNameDescription += ".\nCan not be changed after creation."

	return &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		MaxItems:    1,
		Description: "Identity to apply the rule for.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": &schema.Schema{
					ValidateDiagFunc: validateIdentityType(isCelSupported),
					Type:             schema.TypeString,
					Required:         true,
					Description:      identityTypeDescription,
				},
				"name": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Description: identityNameDescription,
					//ConflictsWith: []string{
					//	"identity.0.group_id",
					//},
				},
				"group_id": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Directory group ID for identity of type GROUP.\nCan not be changed after creation.",
					//ConflictsWith: []string{
					//	"identity.0.name",
					//},
				},
			},
		},
	}
}

func validateIdentityType(isCelSupported bool) func(interface{}, cty.Path) diag.Diagnostics {
	return func(i interface{}, p cty.Path) diag.Diagnostics {

		if i == "USER" ||
			i == "DB_USER" ||
			i == "GROUP" ||
			i == "IDP_GROUP" ||
			i == "DATABRICKS_GROUP" ||
			i == "DATABRICKS_SERVICE_PRINCIPAL" ||
			i == "SNOWFLAKE_ROLE" ||
			i == "EVERYONE" ||
			(isCelSupported && i == "CEL") {
			return diag.Diagnostics{}
		}

		identityTypeErrMsg := fmt.Sprintf("Invalid identity type, `%s`,valid types are: %s", i, validIdentityTypeList(isCelSupported))

		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid value.",
				Detail:   identityTypeErrMsg,
			},
		}
	}
}

func validIdentityTypeList(isCelSupported bool) string {
	identityTypeList := "USER, DB_USER, IDP_GROUP, GROUP, DATABRICKS_GROUP, DATABRICKS_SERVICE_PRINCIPAL, SNOWFLAKE_ROLE"
	if isCelSupported {
		identityTypeList += ", CEL"
	}
	identityTypeList += ", EVERYONE.\nCan not be changed after creation."

	return identityTypeList
}

func resourceDataAccessSecurityPolicies() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "IDs of security policies to apply to this rule. Empty list for default dataset security policies. [ \"none\" ] list for no policies.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func validateAttributeChange(ctx context.Context, old, new, meta interface{}, attrPath string, attrFriendlyPath string) error {
	log.Printf("Validating %s change from %v to %v", attrPath, old, new)

	oldStr, oldOk := old.(string)
	newStr, newOk := new.(string)

	if !oldOk {
		return fmt.Errorf("%s old value is not a string, got: %T", attrFriendlyPath, old)
	}
	if !newOk {
		return fmt.Errorf("%s new value is not a string, got: %T", attrFriendlyPath, new)
	}
	// Only validate if the resource already exists (old value is not empty)
	if oldStr != "" && oldStr != newStr {
		return fmt.Errorf("%s cannot be changed after creation (current: %s, attempted: %s)",
			attrFriendlyPath, oldStr, newStr)
	}
	return nil
}

func validateIdentityTypeChange(ctx context.Context, old, new, meta interface{}) error {
	return validateAttributeChange(ctx, old, new, meta, identityTypePath, "identity.type")
}

func validateIdentityNameChange(ctx context.Context, old, new, meta interface{}) error {
	return validateAttributeChange(ctx, old, new, meta, identityNamePath, "identity.name")
}

func validateIdentityGroupIdChange(ctx context.Context, old, new, meta interface{}) error {
	return validateAttributeChange(ctx, old, new, meta, identityGroupIdPath, "identity.group_id")
}

func resourceToIdentity(resourceIdentity map[string]interface{}) *api.DataAccessIdentity {
	var identity api.DataAccessIdentity

	identity.IdentityType = resourceIdentity["type"].(string)
	identityName := resourceIdentity["name"].(string)
	identityGroupId := resourceIdentity["group_id"].(string)

	if len(identityName) > 0 {
		identity.Identity = identityName
	} else if len(identityGroupId) > 0 {
		identity.Identity = identityGroupId
	}
	return &identity
}

func resourceToDataAccessSecurityPolicies(d *schema.ResourceData) *[]string {
	if raw, ok := d.GetOk("security_policies"); ok {
		in := raw.([]interface{})
		if len(in) > 0 {
			if in[0].(string) == "none" {
				sp := make([]string, 0)
				return &sp
			} else {
				sp := make([]string, len(in))
				for i, v := range in {
					sp[i] = v.(string)
				}
				return &sp
			}
		}
	}
	return nil
}

func dataAccessSecurityPoliciesToResource(securityPolicies *[]string, d *schema.ResourceData) error {
	if securityPolicies != nil {
		if len(*securityPolicies) == 0 {
			sp := []string{"none"}
			return d.Set("security_policies", sp)
		} else {
			return d.Set("security_policies", *securityPolicies)
		}
	} else if _, ok := d.GetOk("security_policies.0"); ok {
		return d.Set("security_policies", []interface{}{})
	}
	return nil
}

func dataAccessUnusedTimeLimitToResource(unusedTimeLimit *api.DataAccessUnusedTimeLimit, d *schema.ResourceData) error {
	if unusedTimeLimit.ShouldRevoke {
		return d.Set("revoke_if_not_used_in_days", unusedTimeLimit.UnusedDaysUntilRevocation)
	}
	return d.Set("revoke_if_not_used_in_days", 0)
}

func dataAccessIdentityToResource(in *api.DataAccessIdentity) *map[string]interface{} {
	out := make(map[string]interface{})
	out["type"] = in.IdentityType

	switch in.IdentityType {
	case "GROUP":
		out["group_id"] = in.Identity
	default: // "IDP_GROUP", "USER", "CEL" and others
		out["name"] = in.Identity
	}
	return &out
}
