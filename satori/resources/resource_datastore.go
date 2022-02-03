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
		Schema: map[string]*schema.Schema{
			"datastore_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Parent ID for dataset permissions.",
			},
			"definition": getDataStoreDefinitionSchema(),
			"access_control_settings": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "DataStore access controls.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_access_control": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Enforce access control to this dataset.",
						},
					},
				},

				//		"enable_user_requests": &schema.Schema{
				//			Type:        schema.TypeBool,
				//			Optional:    true,
				//			Default:     false,
				//			Description: "Allow users to request access to this dataset.",
				//		},
				//		"enable_self_service": &schema.Schema{
				//			Type:        schema.TypeBool,
				//			Optional:    true,
				//			Default:     false,
				//			Description: "Allow users to grant themselves access to this dataset.",
				//		},
				//	},
				//},
			},
			"custom_policy": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "DataStore custom policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     api.CustomPolicyDefaultPriority,
							Description: "DataStore custom policy priority.",
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
				Description: "IDs of security policies to apply to this DataStore.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
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
