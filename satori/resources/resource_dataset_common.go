package resources

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"log"
)

var (
	RelationalLocation          = "relational_location"
	MySqlLocation               = "mysql_location"
	AthenaLocation              = "athena_location"
	MongoLocation               = "mongo_location"
	S3Location                  = "s3_location"
	RelationalLocationType      = "RELATIONAL_LOCATION"
	MySqlLocationType           = "MYSQL_LOCATION"
	AthenaLocationType          = "ATHENA_LOCATION"
	MongoLocationType           = "MONGO_LOCATION"
	S3LocationType              = "S3_LOCATION"
	RelationalTableLocationType = "RELATIONAL_TABLE_LOCATION"
	MySqlTableLocationType      = "MYSQL_TABLE_LOCATION"
	AthenaTableLocationType     = "ATHENA_TABLE_LOCATION"
	MongoTableLocationType      = "MONGO_TABLE_LOCATION"
	S3TableLocationType         = "S3_TABLE_LOCATION"
	Db                          = "db"
	Schema                      = "schema"
	Table                       = "table"
	Catalog                     = "catalog"
	Collection                  = "collection"
	Bucket                      = "bucket"
	ObjectKey                   = "object_key"
	Location                    = "location"
)

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
				"approvers": &schema.Schema{
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Identities of Satori users/groups that will be set as dataset approvers.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"type": &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: ValidateApproverType(false),
								Required:         true,
								Description:      "Approver type, can be either `GROUP` (IdP Group alone) or `USER`",
							},
							"id": &schema.Schema{
								Type:        schema.TypeString,
								Required:    true,
								Description: "The ID of the approver entity",
							},
						},
					},
				},
				"include_location": &schema.Schema{
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Location to include in dataset.",
					Elem:        getDatasetLocationResource(),
				},
				"exclude_location": &schema.Schema{
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Location to exclude from dataset.",
					Elem:        getDatasetLocationResource(),
				},
			},
		},
	}
}

func ValidateApproverType(enableManager bool) func(interface{}, cty.Path) diag.Diagnostics {
	return func(i interface{}, p cty.Path) diag.Diagnostics {

		if i == "USER" || i == "GROUP" || i == "DIRECTORY" || (enableManager && i == "MANAGER") {
			return diag.Diagnostics{}
		}

		message := "Approver type MUST be either 'USER', 'GROUP'"

		if enableManager {
			message = ", 'DIRECTORY' or 'MANAGER'"
		} else {
			message = message + " or 'DIRECTORY'"
		}

		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid value.",
				Detail:   message,
			},
		}
	}
}

func createDataSet(d *schema.ResourceData, c *api.Client) (*api.DataSetOutput, error) {
	input, err := resourceToDataset(d)
	if err != nil {
		return nil, err
	}

	result, err := c.CreateDataSet(input)
	if err != nil {
		return nil, err
	}

	d.SetId(result.Id)

	if err := d.Set("data_policy_id", result.DataPolicyId); err != nil {
		return nil, err
	}

	return result, err
}

func resourceToDataset(d *schema.ResourceData) (*api.DataSet, error) {
	out := api.DataSet{}
	out.Name = d.Get("definition.0.name").(string)

	if v, ok := d.GetOk("definition.0.approvers"); ok {
		out.Approvers = approversInputToResource(v.([]interface{}))
	} else {
		out.Approvers = []api.ApproverIdentity{}
	}

	out.Description = d.Get("definition.0.description").(string)
	if v, ok := d.GetOk("definition.0.owners"); ok {
		owners := v.([]interface{})
		outOwners := make([]string, len(owners))
		for i, owner := range owners {
			if owner == nil {
				return nil, fmt.Errorf("owner can't be empty or null")
			}
			outOwners[i] = owner.(string)
		}
		out.OwnersIds = outOwners
	} else {
		out.OwnersIds = []string{}
	}
	includeLocationOutput, err := resourceToLocations(d, "definition.0.include_location", false)
	if err != nil {
		return nil, err
	}
	out.IncludeLocations = *includeLocationOutput
	excludeLocationOutput, err := resourceToLocations(d, "definition.0.exclude_location", true)
	if err != nil {
		return nil, err
	}
	out.ExcludeLocations = *excludeLocationOutput

	out.CustomPolicy = *resourceToCustomPolicy(d)

	securityPolicies, err := resourceToSecurityPolicies(d)
	if err != nil {
		return nil, err
	}
	out.SecurityPolicies = *securityPolicies

	out.PermissionsEnabled = (*resourceToAccessControl(d)).AccessControlEnabled

	return &out, nil
}

func resourceToLocations(d *schema.ResourceData, mainParamName string, forceLocation bool) (*[]api.DataSetLocation, error) {
	if v, ok := d.GetOk(mainParamName); ok {
		out := make([]api.DataSetLocation, len(v.([]interface{})))
		for i, raw := range v.([]interface{}) {
			inElement := raw.(map[string]interface{})
			outElement, err := resourceToDatasetLocation(inElement, d, forceLocation)
			if err != nil {
				return nil, err
			}
			out[i] = *outElement
		}
		return &out, nil
	}
	out := make([]api.DataSetLocation, 0)
	return &out, nil
}

func resourceToDatasetLocation(inElement map[string]interface{}, d *schema.ResourceData, forceLocation bool) (*api.DataSetLocation, error) {
	outElement := api.DataSetLocation{}
	outElement.DataStoreId = inElement["datastore"].(string)

	err := checkThatOnlyOneLocationFormatExists(inElement, RelationalLocation, Location, forceLocation)
	if err != nil {
		return nil, err
	}

	if len(inElement[RelationalLocation].([]interface{})) > 0 { // deprecated field
		inLocations := inElement[RelationalLocation].([]interface{})
		if len(inLocations) > 0 {
			var location api.DataSetGenericLocation
			err := resourceToGenericLocation(&location, inLocations, RelationalLocationType)
			if err != nil {
				return nil, err
			}
			outElement.Location = &location
		}
	} else if len(inElement[Location].([]interface{})) > 0 { // new field
		inLocations := inElement[Location].([]interface{})
		if len(inLocations) > 0 {
			var location api.DataSetGenericLocation
			err := resourceToLocation(&location, inLocations, false)
			if err != nil {
				return nil, err
			} else {
				outElement.Location = &location
			}
		}
	}
	return &outElement, nil
}

func checkThatOnlyOneLocationFormatExists(inElement map[string]interface{}, deprecatedField string, newField string, forceLocation bool) error {
	if len(inElement[deprecatedField].([]interface{})) > 0 && len(inElement[newField].([]interface{})) > 0 {
		return fmt.Errorf("can not include both fields '%s' and '%s'", deprecatedField, newField)
	}
	if forceLocation && len(inElement[deprecatedField].([]interface{})) == 0 && len(inElement[newField].([]interface{})) == 0 {
		return fmt.Errorf("has to include '%s' field", newField)
	}
	return nil
}

/*
*
The input is for example:
[

	{
	  relational_location: [
	    {
	      db: "db",
	      schema: "schema"
	    }
	  ]
	}

]
*/
func resourceToLocation(location *api.DataSetGenericLocation, locationElem []interface{}, isTableType bool) error {
	inLocationElem := locationElem[0].(map[string]interface{})

	err := checkThatOnlyOneLocationTypeExists(inLocationElem)
	if err != nil {
		return err
	}
	log.Printf("resourceToLocation: %s", inLocationElem)
	if len(inLocationElem[RelationalLocation].([]interface{})) > 0 {
		inLocations := inLocationElem[RelationalLocation].([]interface{})
		if len(inLocations) > 0 {
			locationType := RelationalLocationType
			if isTableType {
				locationType = RelationalTableLocationType
			}
			err := resourceToGenericLocation(location, inLocations, locationType)
			if err != nil {
				return fmt.Errorf("%s in %s", err.Error(), RelationalLocation)
			}
		}
	} else if len(inLocationElem[MySqlLocation].([]interface{})) > 0 {
		inLocations := inLocationElem[MySqlLocation].([]interface{})
		if len(inLocations) > 0 {
			locationType := MySqlLocationType
			if isTableType {
				locationType = MySqlTableLocationType
			}
			err := resourceToGenericLocation(location, inLocations, locationType)
			if err != nil {
				return fmt.Errorf("%s in %s", err.Error(), MySqlLocation)
			}
		}
	} else if len(inLocationElem[AthenaLocation].([]interface{})) > 0 {
		inLocations := inLocationElem[AthenaLocation].([]interface{})
		if len(inLocations) > 0 {
			locationType := AthenaLocationType
			if isTableType {
				locationType = AthenaTableLocationType
			}
			err := resourceToGenericLocation(location, inLocations, locationType)
			if err != nil {
				return fmt.Errorf("%s in %s", err.Error(), AthenaLocation)
			}
		}
	} else if len(inLocationElem[MongoLocation].([]interface{})) > 0 {
		inLocations := inLocationElem[MongoLocation].([]interface{})
		if len(inLocations) > 0 {
			locationType := MongoLocationType
			if isTableType {
				locationType = MongoTableLocationType
			}
			err := resourceToGenericLocation(location, inLocations, locationType)
			if err != nil {
				return fmt.Errorf("%s in %s", err.Error(), MongoLocation)
			}
		}
	} else if len(inLocationElem[S3Location].([]interface{})) > 0 {
		inLocations := inLocationElem[S3Location].([]interface{})
		if len(inLocations) > 0 {
			locationType := S3LocationType
			if isTableType {
				locationType = S3TableLocationType
			}
			err := resourceToGenericLocation(location, inLocations, locationType)
			if err != nil {
				return fmt.Errorf("%s in %s", err.Error(), S3Location)
			}
		}
	}
	return nil
}

func checkThatOnlyOneLocationTypeExists(inLocationElem map[string]interface{}) error {
	countLocationInstances := 0

	if len(inLocationElem[RelationalLocation].([]interface{})) > 0 {
		countLocationInstances += 1
	}
	if len(inLocationElem[MySqlLocation].([]interface{})) > 0 {
		countLocationInstances += 1
	}
	if len(inLocationElem[AthenaLocation].([]interface{})) > 0 {
		countLocationInstances += 1
	}
	if len(inLocationElem[MongoLocation].([]interface{})) > 0 {
		countLocationInstances += 1
	}
	if len(inLocationElem[S3Location].([]interface{})) > 0 {
		countLocationInstances += 1
	}

	if countLocationInstances > 1 {
		return fmt.Errorf("can not include more than one location type from the above: %s, %s, %s, %s, %s", RelationalLocation, MySqlLocation, AthenaLocation, MongoLocation, S3Location)
	}

	return nil
}

func resourceToGenericLocation(location *api.DataSetGenericLocation, inLocations []interface{}, locationType string) error {
	location.Type = locationType

	if inLocations == nil || inLocations[0] == nil {
		log.Printf("inLocations is nil: %s", inLocations)
		return fmt.Errorf("at least one location has contain value")
	}

	inLocation := inLocations[0].(map[string]interface{})
	log.Printf("In location: %s", inLocation)

	if locationType == RelationalLocationType || locationType == RelationalTableLocationType {
		if len(inLocation[Db].(string)) > 0 {
			db := inLocation[Db].(string)
			location.Db = &db
			if len(inLocation[Schema].(string)) > 0 {
				schema1 := inLocation[Schema].(string)
				location.Schema = &schema1
				if len(inLocation[Table].(string)) > 0 {
					table := inLocation[Table].(string)
					location.Table = &table
				}
			}
		}
	} else if locationType == MySqlLocationType || locationType == MySqlTableLocationType {
		if len(inLocation[Db].(string)) > 0 {
			db := inLocation[Db].(string)
			location.Db = &db
			if len(inLocation[Table].(string)) > 0 {
				table := inLocation[Table].(string)
				location.Table = &table
			}
		}
	} else if locationType == AthenaLocationType || locationType == AthenaTableLocationType {
		if len(inLocation[Catalog].(string)) > 0 {
			catalog := inLocation[Catalog].(string)
			location.Catalog = &catalog
			if len(inLocation[Db].(string)) > 0 {
				db := inLocation[Db].(string)
				location.Db = &db
				if len(inLocation[Table].(string)) > 0 {
					table := inLocation[Table].(string)
					location.Table = &table
				}
			}
		}
	} else if locationType == MongoLocationType || locationType == MongoTableLocationType {
		if len(inLocation[Db].(string)) > 0 {
			db := inLocation[Db].(string)
			location.Db = &db
			if len(inLocation[Collection].(string)) > 0 {
				collection := inLocation[Collection].(string)
				location.Collection = &collection
			}
		}
	} else if locationType == S3LocationType || locationType == S3TableLocationType {
		if len(inLocation[Bucket].(string)) > 0 {
			bucket := inLocation[Bucket].(string)
			location.Bucket = &bucket
			if len(inLocation[ObjectKey].(string)) > 0 {
				objectKey := inLocation[ObjectKey].(string)
				location.ObjectKey = &objectKey
			}
		}
	}
	log.Printf("Out location: %s", location)
	return nil
}

func getDataSet(c *api.Client, d *schema.ResourceData) (*api.DataSetOutput, error) {
	result, err, statusCode := c.GetDataSet(d.Id())
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

	definition["approvers"] = approversToResource(&result.Approvers)

	definition["include_location"] = locationsToResource(&result.IncludeLocations, d, "definition.0.include_location", RelationalLocation)
	definition["exclude_location"] = locationsToResource(&result.ExcludeLocations, d, "definition.0.exclude_location", RelationalLocation)

	if err := d.Set("definition", []map[string]interface{}{definition}); err != nil {
		return nil, err
	}
	if err := d.Set("data_policy_id", result.DataPolicyId); err != nil {
		return nil, err
	}

	return result, err
}

func locationsToResource(in *[]api.DataSetLocation, d *schema.ResourceData, prefixFieldName string, deprecatedFieldName string) *[]map[string]interface{} {
	out := make([]map[string]interface{}, len(*in))
	for i, v := range *in {
		outElement := make(map[string]interface{}, 2)
		outElement["datastore"] = v.DataStoreId
		if v.Location != nil {
			// Checks if the state already contains the deprecated field, if so, convert the output to the deprecated format,
			// otherwise convert to the new format
			if _, ok := d.GetOk(fmt.Sprintf("%s.%d.%s", prefixFieldName, i, deprecatedFieldName)); ok { // deprecated field format
				if v.Location != nil && v.Location.Type == RelationalLocationType {
					location := make(map[string]string, 3)
					if v.Location.Db != nil {
						location[Db] = *v.Location.Db
						if v.Location.Schema != nil {
							location[Schema] = *v.Location.Schema
							if v.Location.Table != nil {
								location[Table] = *v.Location.Table
							}
						}
					}
					outElement[RelationalLocation] = []map[string]string{location}
				}
			} else { // new field format
				outElement[Location] = []map[string]interface{}{locationToResource(v.Location)}
			}
		}
		out[i] = outElement
	}
	return &out
}

func locationToResource(genericLocation *api.DataSetGenericLocation) map[string]interface{} {
	locationWrapper := make(map[string]interface{}, 1)
	if genericLocation.Type == RelationalLocationType || genericLocation.Type == RelationalTableLocationType {
		location := make(map[string]string, 3)
		if genericLocation.Db != nil {
			location[Db] = *genericLocation.Db
			if genericLocation.Schema != nil {
				location[Schema] = *genericLocation.Schema
				if genericLocation.Table != nil {
					location[Table] = *genericLocation.Table
				}
			}
		}
		locationWrapper[RelationalLocation] = []map[string]string{location}
	} else if genericLocation.Type == MySqlLocationType || genericLocation.Type == MySqlTableLocationType {
		location := make(map[string]string, 2)
		if genericLocation.Db != nil {
			location[Db] = *genericLocation.Db
			if genericLocation.Table != nil {
				location[Table] = *genericLocation.Table
			}
		}
		locationWrapper[MySqlLocation] = []map[string]string{location}
	} else if genericLocation.Type == AthenaLocationType || genericLocation.Type == AthenaTableLocationType {
		location := make(map[string]string, 3)
		if genericLocation.Catalog != nil {
			location[Catalog] = *genericLocation.Catalog
			if genericLocation.Db != nil {
				location[Db] = *genericLocation.Db
				if genericLocation.Table != nil {
					location[Table] = *genericLocation.Table
				}
			}
		}
		locationWrapper[AthenaLocation] = []map[string]string{location}
	} else if genericLocation.Type == MongoLocationType || genericLocation.Type == MongoTableLocationType {
		location := make(map[string]string, 2)
		if genericLocation.Db != nil {
			location[Db] = *genericLocation.Db
			if genericLocation.Collection != nil {
				location[Collection] = *genericLocation.Collection
			}
		}
		locationWrapper[MongoLocation] = []map[string]string{location}
	} else if genericLocation.Type == S3LocationType || genericLocation.Type == S3TableLocationType {
		location := make(map[string]string, 2)
		if genericLocation.Bucket != nil {
			location[Bucket] = *genericLocation.Bucket
			if genericLocation.ObjectKey != nil {
				location[ObjectKey] = *genericLocation.ObjectKey
			}
		}
		locationWrapper[S3Location] = []map[string]string{location}
	}
	return locationWrapper
}

func updateDataSet(d *schema.ResourceData, c *api.Client) (*api.DataSetOutput, error) {
	input, err := resourceToDataset(d)
	if err != nil {
		return nil, err
	}
	result, err := c.UpdateDataSet(d.Id(), input)
	return result, err
}

func resourceDataSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	if err := c.DeleteDataSet(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
