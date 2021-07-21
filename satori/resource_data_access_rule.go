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
			"parent_data_policy": resourceDataAccessParent(),
			"access_level":       resourceDataAccessLevel(),
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable the rule.",
			},
			"expire_on": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Expire the rule on the given date and time. RFC3339 date format is expected. Time must be in UTC (i.e. YYYY-MM-DD***T***HH:MM:SS***Z***). Empty value = never expire.",
			},
			"revoke_if_not_used_in_days": resourceDataAccessRevokeIfNotUsed(),
			"identity":                   resourceDataAccessIdentity(),
			"security_policies":          resourceDataAccessSecurityPolicies(),
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

	out.Identity = resourceToDataAccessIdentity(d)

	out.SecurityPolicies = resourceToDataAccessSecurityPolicies(d)

	return &out, suspended
}

func resourceDataAccessPermissionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	result, err, statusCode := c.GetDataAccessPermission(d.Id())
	if statusCode == 404 {
		d.SetId("")
		return diags
	}
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
	} else if v, ok := d.GetOk("expire_on"); ok && len(v.(string)) > 0 {
		if err := d.Set("expire_on", nil); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := dataAccessUnusedTimeLimitToResource(&result.UnusedTimeLimit, d); err != nil {
		diag.FromErr(err)
	}

	if err := dataAccessSecurityPoliciesToResource(result.SecurityPolicies, d); err != nil {
		diag.FromErr(err)
	}

	return diags
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
