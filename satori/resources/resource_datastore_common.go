package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

func getDataStoreDataPolicyIdSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Parent ID for DataStore permissions.",
	}
}

func getDataStoreDefinitionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		MaxItems:    1,
		Description: "Parameters for DataStore definition.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "DataStore name.",
				}, "hostname": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Host FQDN name.",
				},
				"port": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Port number description.",
				},
				"owners": &schema.Schema{
					Type:        schema.TypeList,
					Optional:    true,
					Description: "IDs of Satori users that will be set as DataStore owners.",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
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

	if err := d.Set("data_policy_id", result.DataPolicyId); err != nil {
		return nil, err
	}

	return result, err
}

func resourceToDataStore(d *schema.ResourceData) *api.DataStore {
	out := api.DataStore{}
	out.Name = d.Get("definition.0.name").(string)
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

func resourceToDataStoreLocation(inElement map[string]interface{}) api.DataStoreLocation {
	outElement := api.DataStoreLocation{}
	outElement.DataStoreId = inElement["datastore"].(string)
	if inElement["relational_location"] != nil {
		inLocations := inElement["relational_location"].([]interface{})
		if len(inLocations) > 0 {
			//var location api.DataStoreGenericLocation
			//resourceToGenericLocation(&location, inLocations, "RELATIONAL_LOCATION")
			//outElement.Location = &location
		}
	}
	return outElement
}

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

func getDataStore(c *api.Client, d *schema.ResourceData) (*api.DataStoreOutput, error) {
	result, err, statusCode := c.GetDataStore(d.Id())
	if statusCode == 404 {
		d.SetId("")
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	definition := make(map[string]interface{})
	definition["name"] = result.Name
	definition["description"] = result.Description
	definition["owners"] = result.OwnersIds

	//definition["include_location"] = locationsToResource(&result.IncludeLocations)
	//definition["exclude_location"] = locationsToResource(&result.ExcludeLocations)

	if err := d.Set("definition", []map[string]interface{}{definition}); err != nil {
		return nil, err
	}

	if err := d.Set("data_policy_id", result.DataPolicyId); err != nil {
		return nil, err
	}

	return result, err
}

//func locationsToResource(in *[]api.DataStoreLocation) *[]map[string]interface{} {
//  out := make([]map[string]interface{}, len(*in))
//  for i, v := range *in {
//    outElement := make(map[string]interface{}, 2)
//    outElement["datastore"] = v.DataStoreId
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
//    out[i] = outElement
//  }
//  return &out
//}

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
