package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
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

func resourceDataAccessIdentity() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		MaxItems:    1,
		Description: "Identity to apply the rule for.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Identity type, valid types are: USER, IDP_GROUP, GROUP, EVERYONE.\nCan not be changed after creation.",
				},
				"name": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Description: "User/group name for identity types of USER and IDP_GROUP.\nCan not be changed after creation.",
					ConflictsWith: []string{
						"identity.0.group_id",
					},
				},
				"group_id": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Directory group ID for identity of type GROUP.\nCan not be changed after creation.",
					ConflictsWith: []string{
						"identity.0.name",
					},
				},
			},
		},
	}
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

func resourceToDataAccessIdentity(d *schema.ResourceData) *api.DataAccessIdentity {
	var identity api.DataAccessIdentity
	identity.IdentityType = d.Get("identity.0.type").(string)
	if v, ok := d.GetOk("identity.0.name"); ok {
		identity.Identity = v.(string)
	} else if v, ok := d.GetOk("identity.0.group_id"); ok {
		identity.Identity = v.(string)
	} else {
		//for everyone
		identity.Identity = identity.IdentityType
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
	case "IDP_GROUP", "USER":
		out["name"] = in.Identity
	case "GROUP":
		out["group_id"] = in.Identity
	default:
	}
	return &out
}
