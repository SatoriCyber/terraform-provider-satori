package satori

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

func resourceDataSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataSetCreate,
		ReadContext:   resourceDataSetRead,
		UpdateContext: resourceDataSetUpdate,
		DeleteContext: resourceDataSetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Full dataset configuration.",
		Schema: map[string]*schema.Schema{
			"data_policy_id": getDatasetDataPolicyIdSchema(),
			"definition":     getDatasetDefinitionSchema(),
			"access_control_settings": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Dataset access controls.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_access_control": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Enforce access control to this dataset.",
						},
						"enable_user_requests": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Allow users to request access to this dataset.",
						},
						"enable_self_service": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Allow users to grant themselves access to this dataset.",
						},
					},
				},
			},
			"custom_policy": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Dataset custom policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     api.CustomPolicyDefaultPriority,
							Description: "Dataset custom policy priority.",
						},
						"rules_yaml": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom policy rules YAML.",
						},
						"tags_yaml": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom policy tags YAML.",
						},
					},
				},
			},
			"security_policies": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "IDs of security policies to apply to this dataset.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceDataSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	result, err := createDataSet(d, c)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = updateDataPolicy(err, c, result.DataPolicyId, d); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceToCustomPolicy(d *schema.ResourceData) *api.CustomPolicy {
	out := api.CustomPolicy{}
	out.Priority = d.Get("custom_policy.0.priority").(int)
	out.RulesYaml = d.Get("custom_policy.0.rules_yaml").(string)
	out.TagsYaml = d.Get("custom_policy.0.tags_yaml").(string)
	return &out
}

func resourceToAccessControl(d *schema.ResourceData) *api.AccessControl {
	out := api.AccessControl{}
	out.AccessControlEnabled = d.Get("access_control_settings.0.enable_access_control").(bool)
	out.UserRequestsEnabled = d.Get("access_control_settings.0.enable_user_requests").(bool)
	out.SelfServiceEnabled = d.Get("access_control_settings.0.enable_self_service").(bool)
	return &out
}

func resourceToSecurityPolicies(d *schema.ResourceData) *api.SecurityPolicies {
	out := api.SecurityPolicies{}
	if raw, ok := d.GetOk("security_policies"); ok {
		in := raw.([]interface{})
		out.Ids = make([]string, len(in))
		for i, v := range in {
			out.Ids[i] = v.(string)
		}
	}
	return &out
}

func resourceDataSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	result, err := getDataSet(c, d)
	if result == nil && err == nil {
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	resultCustomPolicy, err := c.GetCustomPolicy(result.DataPolicyId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("custom_policy", []map[string]interface{}{*customPolicyToResource(resultCustomPolicy)}); err != nil {
		return diag.FromErr(err)
	}

	resultAccessControl, err := c.GetAccessControl(result.DataPolicyId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("access_control_settings", []map[string]interface{}{*accessControlToResource(resultAccessControl)}); err != nil {
		return diag.FromErr(err)
	}

	resultSecurityPolicies, err := c.GetSecurityPolicies(result.DataPolicyId)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicies := securityPoliciesToResource(resultSecurityPolicies)
	var currentSecurityPoliciesLen = 0
	if v, ok := d.GetOk("security_policies"); ok {
		currentSecurityPoliciesLen = len(v.([]interface{}))
	}
	if !(currentSecurityPoliciesLen == 0 && len(*securityPolicies) == 0) {
		if err := d.Set("security_policies", securityPolicies); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func customPolicyToResource(in *api.CustomPolicy) *map[string]interface{} {
	out := make(map[string]interface{})
	out["priority"] = in.Priority
	out["rules_yaml"] = in.RulesYaml
	out["tags_yaml"] = in.TagsYaml
	return &out
}

func accessControlToResource(in *api.AccessControl) *map[string]interface{} {
	out := make(map[string]interface{})
	out["enable_access_control"] = in.AccessControlEnabled
	out["enable_user_requests"] = in.UserRequestsEnabled
	out["enable_self_service"] = in.SelfServiceEnabled
	return &out
}

func securityPoliciesToResource(in *api.SecurityPolicies) *[]string {
	out := make([]string, len(in.Ids))
	for i, v := range in.Ids {
		out[i] = v
	}
	return &out
}

func resourceDataSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	result, err := updateDataSet(d, c)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = updateDataPolicy(err, c, result.DataPolicyId, d); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateDataPolicy(err error, c *api.Client, id string, d *schema.ResourceData) error {
	if _, err = c.UpdateCustomPolicy(id, resourceToCustomPolicy(d)); err != nil {
		return err
	}

	if _, err = c.UpdateAccessControl(id, resourceToAccessControl(d)); err != nil {
		return err
	}

	if _, err = c.UpdateSecurityPolicies(id, resourceToSecurityPolicies(d)); err != nil {
		return err
	}

	return nil
}
