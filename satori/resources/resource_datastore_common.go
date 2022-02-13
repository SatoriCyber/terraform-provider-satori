package resources

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

var (
	Name                        = "name"
	Hostname                    = "hostname"
	SatoriHostname              = "satori_hostname"
	Id                          = "id"
	ParentId                    = "parent_id"
	IdentityProviderId          = "identity_provider_id"
	DataAccessControllerId      = "dataaccess_controller_id"
	CustomIngressPort           = "custom_ingress_port"
	Port                        = "port"
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
}

func getDataStoreDefinitionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		Id: &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "DataStore name.",
		},
		ParentId: &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "DataStore name.",
		},
		Name: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "DataStore name.",
		}, Hostname: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Host FQDN name.",
		}, DataAccessControllerId: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Host FQDN name.",
		},
		Port: &schema.Schema{
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Port number description.",
		}, CustomIngressPort: &schema.Schema{
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
		}, SatoriHostname: &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IDs of Satori users that will be set as DataStore owners.",
		},
		BaselineSecurityPolicy: {
			Type:        schema.TypeList,
			Optional:    true,
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
						Optional:    true,
						MaxItems:    1,
						Description: "Baseline security policy.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								QueryAction: &schema.Schema{
									Type:        schema.TypeString,
									Default:     "PASS",
									Optional:    true,
									Description: "DataStore custom policy priority.",
								}}}},

					UnsupportedQueriesCategory: {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Baseline security policy.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								QueryAction: &schema.Schema{
									Type:        schema.TypeString,
									Default:     "PASS",
									Optional:    true,
									Description: "DataStore custom policy priority.",
								}}}},

					Exclusions: {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Baseline security policy.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								ExcludedIdentities: &schema.Schema{
									Type:        schema.TypeList,
									Optional:    true,
									Default:     nil,
									Description: "DataStore custom policy priority.",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											IdentityType: &schema.Schema{
												Type:        schema.TypeString,
												Optional:    true,
												Description: "DataStore custom policy priority.",
											}, Identity: &schema.Schema{
												Type:        schema.TypeString,
												Optional:    true,
												Description: "DataStore custom policy priority.",
											},
										}},
								},
								ExcludedQueryPatterns: &schema.Schema{
									Type:        schema.TypeList,
									Optional:    true,
									Default:     nil,
									Description: "DataStore custom policy priority.",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											Pattern: &schema.Schema{
												Type:        schema.TypeString,
												Optional:    true,
												Description: "DataStore custom policy priority.",
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
	input := resourceToDataStore(d)
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
func resourceToDataStore(d *schema.ResourceData) *api.DataStore {
	out := api.DataStore{}
	out.Name = d.Get("name").(string)
	out.Hostname = d.Get("hostname").(string)
	out.SatoriHostname = d.Get(SatoriHostname).(string)
	out.Port = d.Get("port").(int)
	out.CustomIngressPort = d.Get("custom_ingress_port").(int)
	out.IdentityProviderId = d.Get("identity_provider_id").(string)
	out.DataAccessControllerId = d.Get("dataaccess_controller_id").(string)
	out.ProjectIds = convertStringSet(d.Get("project_ids").(*schema.Set))
	re := baselineSecurityPolicyToResource(d.Get("baseline_security_policy").([]interface{}))
	//if re != nil {
	out.BaselineSecurityPolicy = re
	//} else {
	//	out.BaselineSecurityPolicy = nil
	//}
	out.Type = d.Get("type").(string)
	return &out
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

	d.Set("id", result.Id)
	d.Set(Name, result.Name)
	d.Set(Hostname, result.Hostname)
	d.Set(ParentId, result.ParentId)
	d.Set(Type, result.Type)
	d.Set(Port, result.Port)
	d.Set(CustomIngressPort, result.CustomIngressPort)
	d.Set(IdentityProviderId, result.IdentityProviderId)
	d.Set(DataAccessControllerId, result.DataAccessControllerId)
	d.Set(SatoriHostname, result.SatoriHostname)
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

func baselineSecurityPolicyToResource(in []interface{}) *api.BaselineSecurityPolicy {
	var bls api.BaselineSecurityPolicy
	lesa := extractValueFromInterface(in)
	if lesa == nil {
		return nil
	}
	tfMap := deepCopyMap(lesa, true)
	jk, _ := json.Marshal(tfMap)
	err := json.Unmarshal(jk, &bls)
	if err != nil {
		return nil
	}

	return &bls
}

func updateDataStore(d *schema.ResourceData, c *api.Client) (*api.DataStoreOutput, error) {
	input := resourceToDataStore(d)
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
