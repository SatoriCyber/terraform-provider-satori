package resources

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

var (
	Name                        = "name"
	Hostname                    = "hostname"
	Id                          = "id"
	ParentId                    = "parent_id"
	IdentityProviderId          = "identity_provider_id"
	DataAccessControllerId      = "dataaccess_controller_id"
	CustomIngressPort           = "custom_ingress_port"
	OriginPort                  = "origin_port"
	ProjectIds                  = "project_ids"
	BaselineSecurityPolicy      = "baseline_security_policy"
	Type                        = "type"
	UnassociatedQueriesCategory = "unassociated_queries_category"
	UnsupportedQueriesCategory  = "unsupported_queries_category"
	Pattern                     = "pattern"
	ExcludedIdentities          = "excluded_identities"
	Exclusions                  = "exclusions"
	QueryAction                 = "query_action"
	ExcludedQueryPatterns       = "excluded_query_patterns"
	Identity                    = "identity"
	IdentityType                = "identity_type"
)
var TreatAsMap = map[string]bool{
	Exclusions:                  true,
	UnsupportedQueriesCategory:  true,
	UnassociatedQueriesCategory: true,
	BaselineSecurityPolicy:      true,
}

func getDataStoreDefinitionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		Id: &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "DataStore resource id.",
		},
		ParentId: &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Parent resource id.",
		},
		Name: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "DataStore name.",
		}, Hostname: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Data provider's FQDN hostname.", // example: snowflakecomputing.com, xyz.redshift.amazonaws.com:5439/dev
		}, OriginPort: &schema.Schema{
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Port number description.",
		}, DataAccessControllerId: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Host FQDN name.",
		},
		CustomIngressPort: &schema.Schema{
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Port number description.",
		}, IdentityProviderId: &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Port number description.",
		}, ProjectIds: &schema.Schema{
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "ProjectIds list of project IDs",
			Elem:        &schema.Schema{Type: schema.TypeString},
		}, Type: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "IDs of Satori users that will be set as DataStore owners.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		}, IdentityProviderId: &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IDs of Satori users that will be set as DataStore owners.",
		},
		BaselineSecurityPolicy: {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Description: "Baseline security policy.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					Type: &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "DataStore basepolicy .",
						Default:     "BASELINE_POLICY",
					},
					UnassociatedQueriesCategory: {
						Type:        schema.TypeList,
						Required:    true,
						MaxItems:    1,
						Description: "UnassociatedQueriesCategory",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								QueryAction: &schema.Schema{
									Type:        schema.TypeString, //Default:     "PASS",
									Optional:    true,
									Description: "Default policy action for querying locations that are not associated with a dataset.",
								}}}},

					UnsupportedQueriesCategory: {
						Type:        schema.TypeList,
						Required:    true,
						MaxItems:    1,
						Description: "UnsupportedQueriesCategory",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								QueryAction: &schema.Schema{
									Type:        schema.TypeString, //Default:     "PASS",
									Required:    true,
									Description: "Default policy action for unsupported queries and objects",
								}}}},

					Exclusions: {
						Type:        schema.TypeList,
						Required:    true,
						MaxItems:    1,
						Description: "Exempt users and patterns from baseline security policy",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								ExcludedIdentities: &schema.Schema{
									Type:        schema.TypeList,
									Optional:    true,
									Default:     nil,
									Description: "Exempt Users from the Baseline Security Policy\n",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											IdentityType: &schema.Schema{
												Type:        schema.TypeString,
												Optional:    true,
												Description: "USER type are supported",
											}, Identity: &schema.Schema{
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Username",
											},
										}},
								},
								ExcludedQueryPatterns: &schema.Schema{
									Type:        schema.TypeList,
									Optional:    true,
									Default:     nil,
									Description: "Exempt Queries from the Baseline Security Policy",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											Pattern: &schema.Schema{
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Query pattern",
											},
										}},
								},
							},
						}},
				},
			}},
	}
}

func createDataStore(d *schema.ResourceData, c *api.Client) (*api.DataStoreOutput, error) {
	input, err := resourceToDataStore(d)

	if err != nil {
		return nil, err
	}

	result, err := c.CreateDataStore(input)
	if err != nil {
		return nil, err
	}
	d.SetId(result.Id)

	if err := d.Set("id", result.Id); err != nil {
		return nil, err
	}
	return result, err
}

// convert terraform resource defs into datastore type //
func resourceToDataStore(d *schema.ResourceData) (*api.DataStore, error) {

	re, err := baselineSecurityPolicyToResource(d.Get("baseline_security_policy").([]interface{}))
	if err != nil {
		return nil, err
	}

	out := api.DataStore{}
	out.Name = d.Get("name").(string)
	out.Hostname = d.Get("hostname").(string)
	out.OriginPort = d.Get(OriginPort).(int)
	out.CustomIngressPort = d.Get("custom_ingress_port").(int)
	out.IdentityProviderId = d.Get("identity_provider_id").(string)
	out.DataAccessControllerId = d.Get("dataaccess_controller_id").(string)
	out.ProjectIds = convertStringSet(d.Get("project_ids").(*schema.Set))
	out.BaselineSecurityPolicy = re
	//} else {
	//	out.BaselineSecurityPolicy = nil
	//}
	out.Type = d.Get("type").(string)
	return &out, nil
}

// update datastoreoutput
func getDataStore(c *api.Client, d *schema.ResourceData) (*api.DataStoreOutput, error) {
	result, err, statusCode := c.GetDataStore(d.Id())
	if statusCode == 404 {
		d.SetId("")
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	d.Set(Id, result.Id)
	d.Set(Name, result.Name)
	d.Set(Hostname, result.Hostname)
	d.Set(ParentId, result.ParentId)
	d.Set(Type, result.Type)
	d.Set(OriginPort, result.OriginPort)
	d.Set(CustomIngressPort, result.CustomIngressPort)
	d.Set(IdentityProviderId, result.IdentityProviderId)
	d.Set(DataAccessControllerId, result.DataAccessControllerId)
	d.Set(ProjectIds, result.ProjectIds)

	tfMap, err := getBaseLinePolicyOutput(result, err)
	if err != nil {
		return nil, err
	}
	d.Set(BaselineSecurityPolicy, []map[string]interface{}{tfMap})
	return result, err
}

func getBaseLinePolicyOutput(result *api.DataStoreOutput, err error) (map[string]interface{}, error) {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(result.BaselineSecurityPolicy)
	err = json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return nil, err
	}
	tfMap := deepCopyMap(inInterface, false)
	return tfMap, nil
}

func extractValueFromInterface(in []interface{}) map[string]interface{} {
	if len(in) > 0 {
		return in[0].(map[string]interface{})
	} else {
		return nil
	}
}

func baselineSecurityPolicyToResource(in []interface{}) (*api.BaselineSecurityPolicy, error) {
	var bls api.BaselineSecurityPolicy
	lesa := extractValueFromInterface(in)
	if lesa == nil {
		return nil, errors.New("no datastore correct")
	}
	tfMap := deepCopyMap(lesa, true)
	jk, _ := json.Marshal(tfMap)
	err := json.Unmarshal(jk, &bls)
	if err != nil {
		return nil, (err)
	}
	return &bls, nil
}

func updateDataStore(d *schema.ResourceData, c *api.Client) (*api.DataStoreOutput, error) {
	input, err := resourceToDataStore(d)
	if err != nil {
		return nil, err
	}
	result, err := c.UpdateDataStore(d.Id(), input)
	return result, err
}

func resourceDataStoreDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	if err := c.DeleteDataStore(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
