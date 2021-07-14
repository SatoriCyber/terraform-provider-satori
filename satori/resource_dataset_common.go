package satori

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"strings"
)

func getDatasetLocationResource(locationOptional bool) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datastore": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Data store ID.",
			},
			"relational_location": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    locationOptional,
				Required:    !locationOptional,
				MaxItems:    1,
				Description: "Location for a relational data store.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name.",
						},
						"schema": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Schema name.",
						},
						"table": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Table name.",
						},
					},
				},
			},
		},
	}
}

func getDatasetDataPolicyIdSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Parent ID for dataset permissions.",
	}
}

func getDatasetDefinitionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		MaxItems:    1,
		Description: "Parameters for dataset definition.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Dataset name.",
				},
				"description": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Dataset description.",
				},
				"owners": &schema.Schema{
					Type:        schema.TypeList,
					Optional:    true,
					Description: "IDs of Satori users that will be set as dataset owners.",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"include_location": &schema.Schema{
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Location to include in dataset.",
					Elem:        getDatasetLocationResource(true),
				},
				"exclude_location": &schema.Schema{
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Location to exclude from dataset.",
					Elem:        getDatasetLocationResource(false),
				},
			},
		},
	}
}

func createDataSet(d *schema.ResourceData, c *api.Client) (*api.DataSetOutput, error) {
	dataSet := resourceToDataset(d)

	result, err := c.CreateDataSet(dataSet)
	if err != nil {
		return nil, err
	}

	d.SetId(result.Id)

	if err := d.Set("data_policy_id", result.DataPolicyId); err != nil {
		return nil, err
	}

	return result, err
}

func resourceToDataset(d *schema.ResourceData) *api.DataSet {
	out := api.DataSet{}
	out.Name = d.Get("definition.0.name").(string)
	out.Description = d.Get("definition.0.description").(string)
	if v, ok := d.GetOk("definition.0.owners"); ok {
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
			outElement.DataStoreId = inElement["datastore"].(string)
			if inElement["relational_location"] != nil {
				inLocations := inElement["relational_location"].([]interface{})
				if len(inLocations) > 0 {
					inLocation := inLocations[0].(map[string]interface{})
					if len(inLocation["db"].(string)) > 0 {
						outLocation := inLocation["db"].(string)
						if len(inLocation["schema"].(string)) > 0 {
							outLocation += "." + inLocation["schema"].(string)
							if len(inLocation["table"].(string)) > 0 {
								outLocation += "." + inLocation["table"].(string)
							}
						}
						outElement.Location = &outLocation
					}
				}
			}
			out[i] = outElement
		}
		return &out
	}
	out := make([]api.DataStoreLocation, 0)
	return &out
}

func getDataSet(c *api.Client, d *schema.ResourceData) (*api.DataSetOutput, error) {
	result, err := c.GetDataSet(d.Id())
	if err != nil {
		return nil, err
	}

	definition := make(map[string]interface{})
	definition["name"] = result.Name
	definition["description"] = result.Description
	definition["owners"] = result.OwnersIds

	definition["include_location"] = locationsToResource(&result.IncludeLocations)
	definition["exclude_location"] = locationsToResource(&result.ExcludeLocations)

	if err := d.Set("definition", []map[string]interface{}{definition}); err != nil {
		return nil, err
	}

	if err := d.Set("data_policy_id", result.DataPolicyId); err != nil {
		return nil, err
	}

	return result, err
}

func locationsToResource(in *[]api.DataStoreLocation) *[]map[string]interface{} {
	out := make([]map[string]interface{}, len(*in))
	for i, v := range *in {
		outElement := make(map[string]interface{}, 2)
		outElement["datastore"] = v.DataStoreId
		if v.Location != nil {
			parts := strings.Split(*v.Location, ".")
			location := make(map[string]string, 3)
			location["db"] = parts[0]
			if len(parts) > 1 {
				location["schema"] = parts[1]
				if len(parts) > 2 {
					location["table"] = parts[2]
				}
			}
			outElement["relational_location"] = []map[string]string{location}
		}
		out[i] = outElement
	}
	return &out
}

func updateDataSet(d *schema.ResourceData, c *api.Client) (*api.DataSetOutput, error) {
	dataSet := resourceToDataset(d)
	result, err := c.UpdateDataSet(d.Id(), dataSet)
	return result, err
}

func resourceDataSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	err := c.DeleteDataSet(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
