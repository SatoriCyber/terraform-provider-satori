package resources

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"sort"
)

func getDataStoreDataPolicyIdSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Parent ID for DataStore permissions.",
	}
}

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
					},
					"exclusions": &schema.Schema{
						Type:        schema.TypeList,
						Optional:    true,
						Description: "DataStore custom policy priority.",
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
										},
									}}}},
					},
					"unassociated_queries_category": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "DataStore custom policy priority.",
					}, "unsupported_queries_category": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "DataStore custom policy priority.",
					},
				},
			},
		},
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
	out.BaselineSecurityPolicy = baselineSecurityPolicyToResource(d.Get("baseline_security_policy").([]interface{}))
	out.Type = d.Get("type").(string)
	//out.Description = d.Get("definition.0.description").(string)
	//if v, ok := d.GetOk("definition.0.owners"); ok {
	//	owners := v.([]interface{})
	//	outOwners := make([]string, len(owners))
	//	for i, owner := range owners {
	//		outOwners[i] = owner.(string)
	//	}
	//	out.OwnersIds = outOwners
	//} else {
	//	out.OwnersIds = []string{}
	//}
	//
	//out.IncludeLocations = *resourceToLocations(d, "definition.0.include_location")
	//out.ExcludeLocations = *resourceToLocations(d, "definition.0.exclude_location")
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

//func resourceToLocations(d *schema.ResourceData, mainParamName string) *[]api.DataStoreLocation {
//	if v, ok := d.GetOk(mainParamName); ok {
//		out := make([]api.DataStoreLocation, len(v.([]interface{})))
//		for i, raw := range v.([]interface{}) {
//			inElement := raw.(map[string]interface{})
//			outElement := resourceToDataStoreLocation(inElement)
//			out[i] = outElement
//		}
//		return &out
//	}
//	out := make([]api.DataStoreLocation, 0)
//	return &out
//}

//func resourceToGenericLocation(location *api.DataStoreGenericLocation, inLocations []interface{}, locationType string) {
//	location.Type = locationType
//	inLocation := inLocations[0].(map[string]interface{})
//	log.Printf("In location: %s", inLocation)
//
//	if len(inLocation["db"].(string)) > 0 {
//		db := inLocation["db"].(string)
//		location.Db = &db
//		if len(inLocation["schema"].(string)) > 0 {
//			schema := inLocation["schema"].(string)
//			location.Schema = &schema
//			if len(inLocation["table"].(string)) > 0 {
//				table := inLocation["table"].(string)
//				location.Table = &table
//			}
//		}
//	}
//	log.Printf("Out location: %s", location)
//}

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
	d.Set("baseline_security_policy", result.BaselineSecurityPolicy)

	//definition["include_location"] = locationsToResource(&result.IncludeLocations)
	//definition["exclude_location"] = locationsToResource(&result.ExcludeLocations)

	//if err := d.set( definition); err != nil {
	//	return nil, err
	//}

	//if err := d.Set("data_policy_id", result.DataPolicyId); err != nil {
	//	return nil, err
	//}

	return result, err
}

func baselineSecurityPolicyToResource(in []interface{}) api.BaselineSecurityPolicy {
	var bls api.BaselineSecurityPolicy
	var lesa map[string]interface{}
	lesa = in[0].(map[string]interface{})
	bls.Type = lesa["type"].(string)
	fmt.Println(lesa)
	//for i, v := range *in {
	//outElement := new api.BaselineSecurityPolicy()

	//    if v.Location != nil && v.Location.Type == "RELATIONAL_LOCATION" {
	//      location := make(map[string]string, 3)
	//      if v.Location.Db != nil {
	//        location["db"] = *v.Location.Db
	//        if v.Location.Schema != nil {
	//          location["schema"] = *v.Location.Schema
	//          if v.Location.Table != nil {
	//            location["table"] = *v.Location.Table
	//          }
	//        }
	//      }
	//      outElement["relational_location"] = []map[string]string{location}
	//    }
	//out[0] = outElement

	return bls
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
