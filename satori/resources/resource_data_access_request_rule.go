package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

func ResourceDataAccessRequestRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataAccessRequestRuleCreate,
		ReadContext:   resourceDataAccessRequestRuleRead,
		UpdateContext: resourceDataAccessRequestRuleUpdate,
		DeleteContext: resourceDataAccessRequestRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Request rules configuration.",
		Schema: map[string]*schema.Schema{
			"parent_data_policy": resourceDataAccessParent(),
			"access_level":       resourceDataAccessLevel(),
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable the rule.",
			},
			"expire_in": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Rule expiration settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unit_type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unit type for units field, possible values are: MINUTES, HOURS, DAYS, WEEKS, MONTHS, YEARS.",
						},
						"units": &schema.Schema{
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of units of unit_type.",
						},
					},
				},
			},
			"revoke_if_not_used_in_days": resourceDataAccessRevokeIfNotUsed(),
			"identity":                   resourceDataAccessIdentity(false),
			"security_policies":          resourceDataAccessSecurityPolicies(),
			"require_approver_note": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Require from the approver an `approver note` when approving the request ",
				Default:     false,
			},
		},
	}
}

func resourceDataAccessRequestRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	input := resourceToDataAccessRequestRule(d)

	result, err := c.CreateDataAccessRequestRule(d.Get("parent_data_policy").(string), input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*result.Id)

	return diags
}

func resourceToDataAccessRequestRule(d *schema.ResourceData) *api.DataAccessRequestRule {
	out := api.DataAccessRequestRule{}

	out.AccessLevel = d.Get("access_level").(string)
	out.Suspended = !d.Get("enabled").(bool)

	if v, ok := d.GetOk("expire_in.0.units"); ok {
		out.TimeLimit.Units = v.(int)
		out.TimeLimit.UnitType = d.Get("expire_in.0.unit_type").(string)
		out.TimeLimit.ShouldExpire = true
	} else {
		out.TimeLimit.UnitType = "DAYS"
	}

	revokeUnusedIn := d.Get("revoke_if_not_used_in_days").(int)
	if revokeUnusedIn > 0 {
		out.UnusedTimeLimit.ShouldRevoke = true
		out.UnusedTimeLimit.UnusedDaysUntilRevocation = revokeUnusedIn
	}
	resourceIdentity := d.Get("identity.0").(map[string]interface{})
	out.Identity = resourceToIdentity(resourceIdentity)

	out.SecurityPolicies = resourceToDataAccessSecurityPolicies(d)

	out.RequireApproverNote = d.Get("require_approver_note").(bool)

	return &out
}

func resourceDataAccessRequestRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	result, err, statusCode := c.GetDataAccessRequestRule(d.Id())
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
	if err := d.Set("enabled", !result.Suspended); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("parent_data_policy", *result.ParentId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("identity", []map[string]interface{}{*dataAccessIdentityToResource(result.Identity)}); err != nil {
		return diag.FromErr(err)
	}

	if result.TimeLimit.ShouldExpire {
		expireIn := make(map[string]interface{}, 2)
		expireIn["unit_type"] = result.TimeLimit.UnitType
		expireIn["units"] = result.TimeLimit.Units
		if err := d.Set("expire_in", []map[string]interface{}{expireIn}); err != nil {
			return diag.FromErr(err)
		}
	} else if _, ok := d.GetOk("expire_in.0.units"); ok {
		if err := d.Set("expire_in", nil); err != nil {
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

func resourceDataAccessRequestRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	input := resourceToDataAccessRequestRule(d)
	input.Identity = nil //not allowed to be updated
	_, err := c.UpdateDataAccessRequestRule(d.Id(), input)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDataAccessRequestRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	err := c.DeleteDataAccessRequestRule(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
