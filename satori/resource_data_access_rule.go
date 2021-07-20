package satori

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"time"
)

func resourceDataAccessPermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataAccessPermissionCreate,
		ReadContext:   resourceDataAccessPermissionRead,
		UpdateContext: resourceDataAccessPermissionUpdate,
		DeleteContext: resourceDataAccessPermissionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Access rules configuration.",
		Schema: map[string]*schema.Schema{
			"parent_data_policy": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Parent data policy ID, the data_policy_id field of a dataset.",
			},
			"access_level": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access level to grant, valid values are: READ_ONLY, READ_WRITE, OWNER.",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable the rule.",
			},
			"expire_on": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Expire the rule on the given date and time. RFC3339 date format is expected. Time must be in UTC (i.e. YYYY-MM-DD***T***HH:MM:SS***Z***). Empty value = never expire.",
			},
			"revoke_if_not_used_in_days": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Revoke access if rule not used in the last given days. Zero = do not revoke.",
			},
			"identity": &schema.Schema{
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
			},
		},
	}
}

func resourceDataAccessPermissionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	input, suspended := resourceToDataAccessPermission(d)

	result, err := c.CreateDataAccessPermission(d.Get("parent_data_policy").(string), input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*result.Id)

	if suspended != *result.Suspended {
		if _, err := c.UpdateDataAccessPermissionSuspendedStatus(*result.Id, suspended); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceToDataAccessPermission(d *schema.ResourceData) (*api.DataAccessPermission, bool) {
	out := api.DataAccessPermission{}

	out.AccessLevel = d.Get("access_level").(string)
	suspended := !d.Get("enabled").(bool)

	if v, ok := d.GetOk("expire_on"); ok {
		out.TimeLimit.Expiration = &v //on input it is RFC3339 string
		out.TimeLimit.ShouldExpire = true
	}

	revokeUnusedIn := d.Get("revoke_if_not_used_in_days").(int)
	if revokeUnusedIn > 0 {
		out.UnusedTimeLimit.ShouldRevoke = true
		out.UnusedTimeLimit.UnusedDaysUntilRevocation = revokeUnusedIn
	}

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
	out.Identity = &identity

	return &out, suspended
}

func resourceDataAccessPermissionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	result, err := c.GetDataAccessPermission(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("access_level", result.AccessLevel); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled", !*result.Suspended); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("parent_data_policy", *result.ParentId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("identity", []map[string]interface{}{*dataAccessIdentityToResource(result.Identity)}); err != nil {
		return diag.FromErr(err)
	}
	if result.TimeLimit.ShouldExpire && result.TimeLimit.Expiration != nil {
		n := int64((*result.TimeLimit.Expiration).(float64)) //on output it is epoch millis in JS numeric format
		if err := d.Set("expire_on", time.Unix(n/1000, 0).UTC().Format(time.RFC3339)); err != nil {
			return diag.FromErr(err)
		}
	}
	if result.UnusedTimeLimit.ShouldRevoke {
		if err := d.Set("revoke_if_not_used_in_days", result.UnusedTimeLimit.UnusedDaysUntilRevocation); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
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

func resourceDataAccessPermissionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	input, suspended := resourceToDataAccessPermission(d)
	input.Identity = nil //not allowed to be updated
	result, err := c.UpdateDataAccessPermission(d.Id(), input)
	if err != nil {
		return diag.FromErr(err)
	}

	if suspended != *result.Suspended {
		if _, err := c.UpdateDataAccessPermissionSuspendedStatus(*result.Id, suspended); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceDataAccessPermissionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	err := c.DeleteDataAccessPermission(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
