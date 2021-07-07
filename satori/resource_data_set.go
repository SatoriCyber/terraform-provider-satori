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
		Schema: map[string]*schema.Schema{
			"datapolicy_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"definition": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"owners_ids": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"include_location": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"datastore_id": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"location": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"exclude_location": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"datastore_id": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"location": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"access_control_settings": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_access_control": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"enable_user_requests": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"enable_self_service": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"custom_policy": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  api.CustomPolicyDefaultPriority,
						},
						"rules_yaml": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"tags_yaml": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"security_policies": {
				Type:     schema.TypeList,
				Optional: true,
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

	dataSet := resourceToDataset(d)

	result, err := c.CreateDataSet(dataSet)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("datapolicy_id", result.DataPolicyId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.Id)

	_, err = c.UpdateCustomPolicy(result.DataPolicyId, resourceToCustomPolicy(d))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = c.UpdateSecurityPolicies(result.DataPolicyId, resourceToSecurityPolicies(d))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = c.UpdateAccessControl(result.DataPolicyId, resourceToAccessControl(d))
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceToDataset(d *schema.ResourceData) *api.DataSet {
	out := api.DataSet{}
	out.Name = d.Get("definition.0.name").(string)
	out.Description = d.Get("definition.0.description").(string)
	if v, ok := d.GetOk("definition.0.owners_ids"); ok {
		owners := v.([]interface{})
		outOwners := make([]string, len(owners))
		for i, owner := range owners {
			outOwners[i] = owner.(string)
		}
		out.OwnersIds = outOwners
	} else {
		out.OwnersIds = []string{}
	}

	out.IncludeLocations = *resourceToLocations(d, "definition.0.include_location")
	out.ExcludeLocations = *resourceToLocations(d, "definition.0.exclude_location")
	return &out
}

func resourceToLocations(d *schema.ResourceData, mainParamName string) *[]api.DataStoreLocation {
	if v, ok := d.GetOk(mainParamName); ok {
		out := make([]api.DataStoreLocation, len(v.([]interface{})))
		for i, raw := range v.([]interface{}) {
			inElement := raw.(map[string]interface{})
			outElement := api.DataStoreLocation{}
			outElement.DataStoreId = inElement["datastore_id"].(string)
			outElement.Location = inElement["location"].(string)
			out[i] = outElement
		}
		return &out
	}
	out := make([]api.DataStoreLocation, 0)
	return &out
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
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*api.Client)

	result, err := c.GetDataSet(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	definition := make(map[string]interface{})
	definition["name"] = result.Name
	definition["description"] = result.Description
	definition["owners_ids"] = result.OwnersIds

	definition["include_location"] = locationsToResource(&result.IncludeLocations)
	definition["exclude_location"] = locationsToResource(&result.ExcludeLocations)

	if err := d.Set("definition", []map[string]interface{}{definition}); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("datapolicy_id", result.DataPolicyId); err != nil {
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

func locationsToResource(in *[]api.DataStoreLocation) *[]map[string]string {
	out := make([]map[string]string, len(*in))
	for i, v := range *in {
		outElement := make(map[string]string, 2)
		outElement["datastore_id"] = v.DataStoreId
		outElement["location"] = v.Location
		out[i] = outElement
	}
	return &out
}

func resourceDataSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	dataSet := resourceToDataset(d)

	result, err := c.UpdateDataSet(d.Id(), dataSet)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = c.UpdateCustomPolicy(result.DataPolicyId, resourceToCustomPolicy(d))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = c.UpdateAccessControl(result.DataPolicyId, resourceToAccessControl(d))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = c.UpdateSecurityPolicies(result.DataPolicyId, resourceToSecurityPolicies(d))
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDataSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*api.Client)

	err := c.DeleteDataSet(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
