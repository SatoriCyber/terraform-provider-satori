package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func getStringListProp(prop string, d *schema.ResourceData) *[]string {
	if raw, ok := d.GetOk(prop); ok {
		in := raw.([]interface{})
		list := make([]string, len(in))
		for i, v := range in {
			list[i] = v.(string)
		}
		return &list
	}
	list := make([]string, 0)
	return &list
}

func getDatasetLocationResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datastore": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Data store ID.",
			},
			RelationalLocation: &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Deprecated:  "The 'relational_location' field has been deprecated. Please use the 'location' field instead.",
				Description: "Location for a relational data store.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Db: &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name.",
						},
						Schema: &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Schema name.",
						},
						Table: &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Table name.",
						},
					},
				},
			},
			Location: {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Location for a data store.", // TODO loren add description
				Elem:        getLocationResource(),
			},
		},
	}
}

func getLocationResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			RelationalLocation: &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Location for a relational data store.", // TODO loren add description for which data stores it supports
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

func getRelationalLocationResource() *schema.Resource {
	return &schema.Resource{
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
	}
}
