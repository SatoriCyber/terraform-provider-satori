package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

func ResourceDataSet() *schema.Resource {
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
					},
				},
			},
			"custom_policy": {
				Type:        schema.TypeList,
				Optional:    true,
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

	_, err := createDataSet(d, c)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceToCustomPolicy(d *schema.ResourceData) *api.CustomPolicy {
	out := api.CustomPolicy{}
	priority := d.Get("custom_policy.0.priority").(int)
	_, ok := d.GetOk("custom_policy")
	if priority == 0 && !ok {
		out.Priority = api.CustomPolicyDefaultPriority
	} else {
		out.Priority = priority
	}
	out.RulesYaml = d.Get("custom_policy.0.rules_yaml").(string)
	out.TagsYaml = d.Get("custom_policy.0.tags_yaml").(string)
	return &out
}

func resourceToAccessControl(d *schema.ResourceData) *api.AccessControl {
	out := api.AccessControl{}
	out.AccessControlEnabled = d.Get("access_control_settings.0.enable_access_control").(bool)
	return &out
}

func resourceToSecurityPolicies(d *schema.ResourceData) *api.SecurityPolicies {
	out := api.SecurityPolicies{}
	out.Ids = *getStringListProp("security_policies", d)
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

	resultCustomPolicy := result.CustomPolicy

	if !(checkIfDefaultCustomPolicy(&resultCustomPolicy) && checkIfStateCustomPolicyIsDefault(d)) {
		if err := d.Set("custom_policy", []map[string]interface{}{*customPolicyToResource(&resultCustomPolicy)}); err != nil {
			return diag.FromErr(err)
		}
	}

	resultAccessControl := api.AccessControl{}
	resultAccessControl.AccessControlEnabled = result.PermissionsEnabled

	if err := d.Set("access_control_settings", []map[string]interface{}{*accessControlToResource(&resultAccessControl)}); err != nil {
		return diag.FromErr(err)
	}

	resultSecurityPolicies := result.SecurityPolicies
	if err := setStringListProp(&resultSecurityPolicies.Ids, "security_policies", d); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func checkIfStateCustomPolicyIsDefault(d *schema.ResourceData) bool {
	_, ok := d.GetOk("custom_policy")
	return !ok || checkIfDefaultCustomPolicy(resourceToCustomPolicy(d))
}

func checkIfDefaultCustomPolicy(policy *api.CustomPolicy) bool {
	return policy == nil || (policy.Priority == api.CustomPolicyDefaultPriority && policy.RulesYaml == "" && policy.TagsYaml == "")
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
	return &out
}

func resourceDataSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	_, err := updateDataSet(d, c)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
