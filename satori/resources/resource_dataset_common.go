package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"log"
	"strings"
)

var (
	RelationalLocation          = "relational_location"       // deprecated
	MySqlLocation               = "mysql_location"            // deprecated
	AthenaLocation              = "athena_location"           // deprecated
	MongoLocation               = "mongo_location"            // deprecated
	S3Location                  = "s3_location"               // deprecated
	RelationalLocationType      = "RELATIONAL_LOCATION"       // deprecated
	MySqlLocationType           = "MYSQL_LOCATION"            // deprecated
	AthenaLocationType          = "ATHENA_LOCATION"           // deprecated
	MongoLocationType           = "MONGO_LOCATION"            // deprecated
	S3LocationType              = "S3_LOCATION"               // deprecated
	RelationalTableLocationType = "RELATIONAL_TABLE_LOCATION" // deprecated
	MySqlTableLocationType      = "MYSQL_TABLE_LOCATION"      // deprecated
	AthenaTableLocationType     = "ATHENA_TABLE_LOCATION"     // deprecated
	MongoTableLocationType      = "MONGO_TABLE_LOCATION"      // deprecated
	S3TableLocationType         = "S3_TABLE_LOCATION"         // deprecated
	Db                          = "db"                        // deprecated
	Schema                      = "schema"                    // deprecated
	Table                       = "table"                     // deprecated
	Catalog                     = "catalog"                   // deprecated
	Collection                  = "collection"                // deprecated
	Bucket                      = "bucket"                    // deprecated
	ObjectKey                   = "object_key"                // deprecated
	Location                    = "location"                  // deprecated
	LocationPath                = "location_path"
	LocationParts               = "location_parts"
	LocationPartsFull           = "location_parts_full"
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

	err := checkThatOnlyOneLocationFormatExists(inElement, Location, LocationPath, LocationParts, LocationPartsFull, forceLocation)
	if err != nil {
		return nil, err
	}

	inElementPrint, _ := json.Marshal(inElement)
	log.Printf("The inElement presentation `%s`, length: %s", inElementPrint)

	if len(inElement[Location].([]interface{})) > 0 { // deprecated field
		inLocations := inElement[Location].([]interface{})
		if len(inLocations) > 0 {
			var location api.DataSetGenericLocation
			err := resourceToLocation(&location, inLocations, false)
			if err != nil {
				return nil, err
			}
			outElement.Location = &location
		}
	} else if len(inElement[LocationParts].([]interface{})) > 0 { // new field for LocationParts
		inLocations := inElement[LocationParts].([]interface{})
		if len(inLocations) > 0 {
			var location []api.LocationPath
			err := resourcePartsToLocationPath(&location, inLocations, false)
			if err != nil {
				return nil, err
			}
			outElement.LocationPath = location
		}
	} else if len(inElement[LocationPartsFull].([]interface{})) > 0 { // new field for LocationPartsFull
		inLocations := inElement[LocationPartsFull].([]interface{})
		if len(inLocations) > 0 {
			var location []api.LocationPath
			err := resourcePartsFullToLocationPath(&location, inLocations, false)
			if err != nil {
				return nil, err
			}
			outElement.LocationPath = location
		}
	} else if inElement[LocationPath] != nil && len(inElement[LocationPath].(string)) > 0 { // new string field, ignore if empty
		// terraform value will be always not nil, so we need to check the length and consider it as does not exist if empty
		inLocationStr := inElement[LocationPath].(string)

		var location []api.LocationPath
		log.Printf("found %s location path with length %d", inLocationStr, len(inLocationStr))
		err := resourceStrToLocationPath(&location, inLocationStr, false)
		if err != nil {
			return nil, err
		}
		outElement.LocationPath = location
	}
	return &outElement, nil
}

func checkThatOnlyOneLocationFormatExists(inElement map[string]interface{},
	deprecatedField string,
	newField1 string,
	newField2 string,
	newField3 string,
	forceLocation bool) error {

	log.Printf("checkThatOnlyOneLocationFormatExists started.")

	hasDeprecatedField := inElement[deprecatedField] != nil && len(inElement[deprecatedField].([]interface{})) > 0
	// this is a string field
	hasNewField1 := inElement[newField1] != nil && len(inElement[newField1].(string)) > 0
	hasNewField2 := inElement[newField2] != nil && len(inElement[newField2].([]interface{})) > 0
	hasNewField3 := inElement[newField3] != nil && len(inElement[newField3].([]interface{})) > 0

	trueCount := 0
	if hasDeprecatedField {
		trueCount++
	}
	if hasNewField1 {
		trueCount++
	}
	if hasNewField2 {
		trueCount++
	}
	if hasNewField3 {
		trueCount++
	}

	if trueCount > 2 {
		return fmt.Errorf("can not include more than 1 field of '%s', '%s', '%s' or '%s'", deprecatedField, newField1, newField3, newField3)
	}

	if forceLocation && len(inElement[deprecatedField].([]interface{})) == 0 && len(inElement[newField1].([]interface{})) == 0 && len(inElement[newField2].([]interface{})) == 0 && len(inElement[newField3].([]interface{})) == 0 {
		return fmt.Errorf("has to include '%s' or '%s' or '%s' field", newField1, newField2, newField3)
	}
	log.Printf("checkThatOnlyOneLocationFormatExists ended.")
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

func resourceStrToLocationPath(location *[]api.LocationPath, locationElem string, isTableType bool) error {
	log.Printf("resourcePartsToLocationPath: got %s", locationElem)
	path := strings.Split(locationElem, `.`)
	for _, element := range path {
		locationPath := api.LocationPath{Name: element}
		*location = append(*location, locationPath)
	}
	return nil
}

func resourcePartsToLocationPath(location *[]api.LocationPath, locationElem []interface{}, isTableType bool) error {
	log.Printf("resourcePartsToLocationPath: got %s", locationElem)
	for _, element := range locationElem {
		locationPath := api.LocationPath{Name: element.(string)}
		*location = append(*location, locationPath)
	}
	log.Printf("resourcePartsToLocationPath: returned %s", LocationPath)
	return nil
}

func resourcePartsFullToLocationPath(location *[]api.LocationPath, locationElem []interface{}, isTableType bool) error {
	log.Printf("resourcePartsFullToLocationPath: got locationElem: %s", locationElem)
	for _, element := range locationElem {
		log.Printf("resourcePartsFullToLocationPath: element: %s", element)
		part := element.(map[string]interface{})
		log.Printf("resourcePartsFullToLocationPath: part: %s", part)
		if len(part) > 0 {
			name := part["name"].(string)
			typeStr := part["type"].(string)
			log.Printf("resourcePartsFullToLocationPath: part name: %s", name)
			log.Printf("resourcePartsFullToLocationPath: part type: %s", typeStr)
			//locationPath := api.LocationPath{Name: name}
			locationPath := api.LocationPath{}
			locationPath.Name = name
			locationPath.Type = typeStr
			log.Printf("resourcePartsFullToLocationPath: locationPath: %s", locationPath)
			*location = append(*location, locationPath)
		}
	}
	log.Printf("Out location: %s", location)
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

	definition["include_location"] = locationsToResource(&result.IncludeLocations, d, "definition.0.include_location", Location)
	definition["exclude_location"] = locationsToResource(&result.ExcludeLocations, d, "definition.0.exclude_location", Location)

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
			log.Printf("locationsToResource: for old deprecated format, %s.%d.%s - %s", prefixFieldName, i, deprecatedFieldName, v.Location)
			if _, ok := d.GetOk(fmt.Sprintf("%s.%d.%s", prefixFieldName, i, deprecatedFieldName)); ok { // deprecated field format
				outElement[Location] = []map[string]interface{}{locationToResource(v.Location)}
			}
		} else if v.LocationPath != nil { // new field format
			if configuredLocationPath, ok := d.GetOk(fmt.Sprintf("%s.%d.%s", prefixFieldName, i, LocationPath)); ok { // new LocationPath was configured
				log.Printf("locationsToResource: new format for %s was found: %s", LocationPath, configuredLocationPath)
				outElement[LocationPath] = locationPathToLocationPathResource(v.LocationPath)
			} else if configuredLocationPath, ok := d.GetOk(fmt.Sprintf("%s.%d.%s", prefixFieldName, i, LocationParts)); ok { // new LocationPath was configured
				log.Printf("locationsToResource: new format for %s was found: %s", LocationParts, configuredLocationPath)
				outElement[LocationParts] = locationPathToLocationPartsResource(v.LocationPath)
			} else if configuredLocationPath, ok := d.GetOk(fmt.Sprintf("%s.%d.%s", prefixFieldName, i, LocationPartsFull)); ok { // new LocationPath was configured
				log.Printf("locationsToResource: new format for %s was found: %s", LocationPartsFull, configuredLocationPath)
				outElement[LocationPartsFull] = locationPathToLocationPartsFullResource(v.LocationPath)
			} else {
				log.Printf("got an unknown format for locationPath")
			}

		}
		out[i] = outElement
	}
	return &out
}

func locationToResource(genericLocation *api.DataSetGenericLocation) map[string]interface{} {
	log.Printf("locationToResource started")

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

func locationPathToLocationPathResource(genericLocation []api.LocationPath) string {
	log.Printf("locationPathToLocationPathResource started")

	locationParts := []string{}
	for _, u := range genericLocation {
		locationParts = append(locationParts, u.Name)
	}
	locationJoinedPath := strings.Join(locationParts, ".")
	log.Printf("locationPathToLocationPathResource found %s: ", locationJoinedPath)
	return locationJoinedPath
}

func locationPathToLocationPartsResource(genericLocation []api.LocationPath) []string {
	log.Printf("locationPathToLocationPartsResource started")
	locationParts := []string{}
	for _, u := range genericLocation {
		locationParts = append(locationParts, u.Name)
	}
	log.Printf("locationPathToLocationPartsResource found %s: ", locationParts)
	return locationParts
}

func locationPathToLocationPartsFullResource(genericLocation []api.LocationPath) []interface{} {
	log.Printf("locationPathToLocationPartsFullResource started")

	elementNumber := len(genericLocation)
	log.Printf("locationPathToLocationPartsFullResource, fournd %d elements", elementNumber)

	locationParts := make([]interface{}, elementNumber)
	for i, u := range genericLocation {
		locationParts[i] = make(map[string]interface{}, 2)
		locationParts[i].(map[string]interface{})[Name] = u.Name
		locationParts[i].(map[string]interface{})[Type] = u.Type
	}
	log.Printf("locationPathToLocationPartsFullResource found %s: ", locationParts)
	return locationParts
}

func updateDataSet(d *schema.ResourceData, c *api.Client) (*api.DataSetOutput, error) {
	input, err := resourceToDataset(d)
	inputJson, err := json.Marshal(input)
	log.Printf("updateDataSet: %s", inputJson)
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
