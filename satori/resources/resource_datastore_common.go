package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

func getDataStoreDefinitionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "DataStore name.",
		},
		"parentid": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "DataStore name.",
		},
		"name": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "DataStore name.",
		}, "hostname": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Host FQDN name.",
		}, "dataaccess_controller_id": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Host FQDN name.",
		},
		"port": &schema.Schema{
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Port number description.",
		}, "custom_ingress_port": &schema.Schema{
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Port number description.",
		}, "identity_provider_id": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Port number description.",
		}, "project_ids": &schema.Schema{
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "ProjectIds list of project IDs",
			Elem:        &schema.Schema{Type: schema.TypeString},
		}, "tags": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IDs of Satori users that will be set as DataStore owners.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		}, "type": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "IDs of Satori users that will be set as DataStore owners.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		}, "rules": &schema.Schema{
			Type:        schema.TypeList,
			Optional:    true,
			Description: "IDs of Satori users that will be set as DataStore owners.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		}, "identityproviderid": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IDs of Satori users that will be set as DataStore owners.",
		},
		//"include_location": &schema.Schema{
		//	Type:        schema.TypeList,
		//	Optional:    true,
		//	Description: "Location to include in DataStore.",
		//	Elem:        getDataStoreLocationResource(true),
		//},
		//"exclude_location": &schema.Schema{
		//	Type:        schema.TypeList,
		//	Optional:    true,
		//	Description: "Location to exclude from DataStore.",
		//	Elem:        getDataStoreLocationResource(false),
		//},
		"baseline_security_policy": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Description: "Baseline security policy.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "DataStore basepolicy .",
						Default:     "BASELINE-POLICY",
					},
					"unassociated_queries_category": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Baseline security policy.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"query_action": &schema.Schema{
									Type:        schema.TypeString,
									Optional:    true,
									Description: "DataStore custom policy priority.",
								}}}},

					"unsupported_queries_category": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Baseline security policy.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"query_action": &schema.Schema{
									Type:        schema.TypeString,
									Optional:    true,
									Description: "DataStore custom policy priority.",
								}}}},

					"exclusions": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Baseline security policy.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"excluded_identities": &schema.Schema{
									Type:        schema.TypeList,
									Optional:    true,
									Description: "DataStore custom policy priority.",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"identity_type": &schema.Schema{
												Type:        schema.TypeString,
												Optional:    true,
												Description: "DataStore custom policy priority.",
											}, "identity": &schema.Schema{
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
	out.Port = d.Get("port").(int)
	out.CustomIngressPort = d.Get("custom_ingress_port").(int)
	out.IdentityProviderId = d.Get("identity_provider_id").(string)
	out.DataAccessControllerId = d.Get("dataaccess_controller_id").(string)
	//out.ProjectIds, ok = d.GetOk("definition.0.project_ids")
	//ok{}
	out.ProjectIds = convertStringSet(d.Get("project_ids").(*schema.Set))
	re := baselineSecurityPolicyToResource(d.Get("baseline_security_policy").([]interface{}))
	if re != nil {
		out.BaselineSecurityPolicy = *re
	}
	out.Type = d.Get("type").(string)

	return &out
}
func convertStringSet(set *schema.Set) []string {
	s := make([]string, 0, set.Len())
	for _, v := range set.List() {
		s = append(s, v.(string))
	}
	sort.Strings(s)

	return s
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
	/// update output from request
	//definition := make(map[string]interface{})
	d.Set("id", result.Id)
	d.Set("name", result.Name)
	d.Set("hostname", result.Hostname)
	d.Set("parentid", result.ParentId)
	d.Set("type", result.Type)
	d.Set("port", result.Port)
	d.Set("custom_ingress_port", result.CustomIngressPort)
	d.Set("identity_provider_id", result.IdentityProviderId)
	d.Set("dataaccess_controller_id", result.DataAccessControllerId)
	d.Set("project_ids", result.ProjectIds)

	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(result.BaselineSecurityPolicy)
	err = json.Unmarshal(inrec, &inInterface)

	//basepolicy := []map[string]interface{}{{
	//	"type":                          result.BaselineSecurityPolicy.Type,
	//	"unsupported_queries_category":  []map[string]interface{}{{"query_action": result.BaselineSecurityPolicy.UnsupportedQueriesCategory.QueryAction}},
	//	"unassociated_queries_category": []map[string]interface{}{{"query_action": result.BaselineSecurityPolicy.UnassociatedQueriesCategory.QueryAction}},
	//}}
	asea := CopyMap(inInterface)

	d.Set("baseline_security_policy", []map[string]interface{}{asea})
	//d.Set("baseline_security_policy", basepolicy)

	return result, err
}
func resNameTfConver(in string) string {
	var tfRegExp = `([A-Z])`
	var re = regexp.MustCompile(tfRegExp)
	s := strings.ToLower(string(re.ReplaceAll([]byte(in), []byte(`_$1`))))
	return (s)
}

func CopyMap(m map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[resNameTfConver(k)] = []map[string]interface{}{CopyMap(vm)}
		} else {
			cp[resNameTfConver(k)] = v
		}
	}

	return cp
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
	bls.Type = lesa["type"].(string)
	if reflect.ValueOf(lesa["unassociated_queries_category"]).Len() > 0 {
		var uaqc api.UnassociatedQueriesCategory
		query_action := (lesa["unassociated_queries_category"]).([]interface{})
		uaqc.QueryAction = extractValueFromInterface(query_action)["query_action"].(string)
		bls.UnassociatedQueriesCategory = uaqc
	}
	if reflect.ValueOf(lesa["unsupported_queries_category"]).Len() > 0 { //	bls.UnsupportedQueriesCategory = lesa["unassociated_queries_category"].(api.UnsupportedQueriesCategory)
		var uaqc api.UnsupportedQueriesCategory
		query_action := (lesa["unsupported_queries_category"]).([]interface{})
		uaqc.QueryAction = extractValueFromInterface(query_action)["query_action"].(string)
		bls.UnsupportedQueriesCategory = uaqc
	}
	if lesa["exclusions"] != nil { //	bls.UnsupportedQueriesCategory = lesa["unassociated_queries_category"].(api.UnsupportedQueriesCategory)
		//var uaqc api.Exclusions
		exclusions := lesa["exclusions"].([]interface{})
		i := extractValueFromInterface(exclusions)["excluded_identities"].([]interface{})
		fmt.Println(i)
		//uaqc.ExcludedIdentities = nil
		var tempIden api.ExcludedIdentities
		tempIden.Identity = "user"
		tempIden.IdentityType = "USER"
		bls.Exclusions.ExcludedIdentities = []api.ExcludedIdentities{tempIden} //i
	}
	fmt.Println(lesa)
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
