package resources

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"log"
	"regexp"
	"sort"
	"strings"
)

// converting API objects back&forth to TF objects
// the function resolves the limitation that TypeMap can't have the TypeList elements
// !!! tf resources and api structures have to be similar !!!
func biTfApiConverter(m map[string]interface{}, camelCase bool) map[string]interface{} {
	currentMap := make(map[string]interface{})
	for k, v := range m {
		vm, okMapInterface := v.(map[string]interface{})
		arrNullInterface, okArrNullInterface := v.([]interface{})
		if (v) == nil && !okArrNullInterface && !okMapInterface {
			currentMap[resNameTfConvert(k, camelCase)] = nil
		} else if okMapInterface {
			currentMap[resNameTfConvert(k, camelCase)] = []map[string]interface{}{biTfApiConverter(vm, camelCase)}
		} else if okArrNullInterface {
			var mapFromNullInterface []map[string]interface{}
			for _, curNullInterface := range arrNullInterface {
				if curNullInterface != nil {
					if currVal := biTfApiConverter(curNullInterface.(map[string]interface{}), camelCase); currVal != nil {
						mapFromNullInterface = append(mapFromNullInterface, currVal)
					}
				}
			}
			if !TreatAsMap[k] {
				currentMap[resNameTfConvert(k, camelCase)] = mapFromNullInterface
			} else {
				var newMapArrayInterface []interface{}
				for _, curRecord := range mapFromNullInterface {
					for _, vaa := range curRecord {
						newMapArrayInterface = append(newMapArrayInterface, vaa)
					}
				}
				if len(mapFromNullInterface) != 0 {
					currentMap[resNameTfConvert(k, camelCase)] = mapFromNullInterface[0]
				} else {
					currentMap[resNameTfConvert(k, camelCase)] = map[string]interface{}{}
				}
			}
		} else {
			currentMap[resNameTfConvert(k, camelCase)] = v
		}
	}
	return currentMap
}

// converts name from camelCase to tf underscore style
func resNameTfConvert(in string, camelCase bool) string {
	if camelCase == true {
		return convertToCamelCase(in)
	} else {
		var tfRegExp = `([A-Z])`
		var re = regexp.MustCompile(tfRegExp)
		s := strings.ToLower(string(re.ReplaceAll([]byte(in), []byte(`_$1`))))
		return s
	}
}
func convertToCamelCase(myString string) string {
	var tfRegExp = `(_)([a-z])`
	re := regexp.MustCompile(tfRegExp).ReplaceAllStringFunc(myString, strings.ToUpper)
	return strings.Replace(re, "_", "", -1)
}

// convert terraform set of strings to string array
func convertStringSet(set *schema.Set) []string {
	s := make([]string, 0, set.Len())
	for _, v := range set.List() {
		s = append(s, v.(string))
	}
	sort.Strings(s)
	return s
}

func setNullableStringProp(in *string, prop string, d *schema.ResourceData) error {
	if in != nil {
		if err := d.Set(prop, in); err != nil {
			return err
		}
	} else {
		if v, ok := d.GetOk(prop); ok && len(v.(string)) > 0 {
			if err := d.Set(prop, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

func setStringListProp(in *[]string, prop string, d *schema.ResourceData) error {
	var list []string
	if in != nil {
		list = *in
	}

	var currentLen = 0
	if v, ok := d.GetOk(prop); ok {
		currentLen = len(v.([]interface{}))
	}
	if !(currentLen == 0 && len(list) == 0) {
		return d.Set(prop, &list)
	}

	return nil
}

func setMapProp(in *map[string]interface{}, prop string, d *schema.ResourceData) error {
	if len(*in) > 0 {
		if err := d.Set(prop, []map[string]interface{}{*in}); err != nil {
			return err
		}
	} else {
		if err := d.Set(prop, nil); err != nil {
			return err
		}
	}
	return nil
}

func getStringListProp(prop string, d *schema.ResourceData) (*[]string, error) {
	if raw, ok := d.GetOk(prop); ok {
		in := raw.([]interface{})
		list := make([]string, len(in))
		for i, v := range in {
			log.Printf("getStringListProp, v=: %s", v)
			if v == nil {
				return nil, fmt.Errorf("can't be empty")
			}
			list[i] = v.(string)
		}
		return &list, nil
	}
	list := make([]string, 0)
	return &list, nil
}

func getDatasetLocationResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datastore": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Data store ID.",
			},
			Location: {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				MinItems:    1,
				Deprecated:  "The 'location' field has been deprecated. Please use the 'location_path', `location_parts` or `location_parts_full` fields instead.",
				Description: "Location for a data store. Can include only one location type field from the above: relational_location, mysql_location, athena_location, mongo_location and s3_location . Conflicts with 'location_path', 'location_parts' and 'location_parts_full' fields.",
				Elem:        getLocationResource(),
			},
			LocationPath: {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The short presentation of the location path in the data store. Includes `.` separated string when part types are defined with default definitions. For example 'a.b.c' in Snowflake data store will path to table 'a' under schema 'b' under database 'a'.  Conflicts with 'location', 'location_parts', and 'location_parts_full' fields.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			LocationParts: {
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    1,
				Description: "The part separated location path in the data store. Includes an array of path parts when part types are defined with default definitions. For example ['a', 'b', 'c'] in Snowflake data store will path to table 'a' under schema 'b' under database 'a'. Conflicts with 'location', 'location_path', and 'location_parts_full' fields",
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: StringIsNotWhiteSpaceInArray,
				},
			},
			LocationPartsFull: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The full location path definition in the data store. Includes an array of objects with path name and path type. Can be used when the path type should be defined explicitly and not as defined by default. For example [{name= 'a', type= 'DATABASE'},{name= 'b', type= 'SCHEMA'},{name= 'view.c', type= 'VIEW'}]. Conflicts with 'location', 'location_path', and 'location_parts' fields.",
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "The name of the location part.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
						},
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "The type of the location part. Optional values: TABLE, COLUMN, SEMANTIC_MODEL, REPORT, DASHBOARD, DATABASE, SCHEMA, JSON_PATH, WAREHOUSE, ENDPOINT, TYPE, FIELD, EXTERNAL_LOCATION, CATALOG, BUCKET, OBJECT, COLLECTION, VIEW, etc",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
						},
					},
				},
			},
		},
	}
}

//func StringIsNotWhiteSpaceInArray(v interface{}, p cty.Path) diag.Diagnostics {
//	value, ok := v.(string)
//
//	var diags diag.Diagnostics
//
//	if !ok {
//		diag := diag.Diagnostic{
//			Severity:      diag.Error,
//			Summary:       "Wrong value type",
//			Detail:        fmt.Sprintf("expected type of %s to be string", p),
//			AttributePath: p,
//		}
//		diags = append(diags, diag)
//	} else {
//
//		if strings.TrimSpace(value) == "" {
//			errorMessage := fmt.Sprintf("value is expected to not be an empty string or whitespace.")
//
//			attr, okIndexStepCasting := p[len(p)-1].(cty.IndexStep)
//			if okIndexStepCasting && attr.Key.AsBigFloat() != nil {
//				index := attr.Key.AsBigFloat().String()
//				errorMessage = fmt.Sprintf("value at index %s is expected to not be an empty string or whitespace.", index)
//			}
//
//			diag := diag.Diagnostic{
//				Severity:      diag.Error,
//				Summary:       "Empty string is not allowed",
//				Detail:        errorMessage,
//				AttributePath: p,
//			}
//			diags = append(diags, diag)
//		}
//	}
//
//	return diags
//}

func getLocationResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			RelationalLocation: &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Location for a relational data store.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Db: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name.",
						},
						Schema: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Schema name.",
						},
						Table: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Table name.",
						},
					},
				},
			},
			MySqlLocation: {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Location for MySql and MariaDB data stores.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Db: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name.",
						},
						Table: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Table name.",
						},
					},
				},
			},
			AthenaLocation: {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Location for Athena data store.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Catalog: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Catalog name.",
						},
						Db: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database name.",
						},
						Table: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Table name.",
						},
					},
				},
			},
			MongoLocation: {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Location for MongoDB data store.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Db: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name.",
						},
						Collection: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Collection name.",
						},
					},
				},
			},
			S3Location: {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Location for S3 data store.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Bucket: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Bucket name.",
						},
						ObjectKey: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Object Key name.",
						},
					},
				},
			},
		},
	}
}

func approversInputToResource(approvers []interface{}) []api.ApproverIdentity {
	approversOutput := make([]api.ApproverIdentity, len(approvers))
	for i, approver := range approvers {
		tmp := approver.(map[string]interface{})
		var mappedApprover api.ApproverIdentity
		mappedApprover.Id = tmp["id"].(string)
		if mappedApprover.Id != "MANAGER" {
			mappedApprover.Type = tmp["type"].(string)
		}
		approversOutput[i] = mappedApprover
	}
	return approversOutput
}

func approversToResource(approvers *[]api.ApproverIdentity) []interface{} {
	mappedApprovers := make([]interface{}, len(*approvers))

	for i, approver := range *approvers {
		approverMap := make(map[string]string)
		approverMap["type"] = approver.Type
		if approver.Type != "MANAGER" {
			approverMap["id"] = approver.Id
		}
		mappedApprovers[i] = approverMap
	}

	return mappedApprovers
}
