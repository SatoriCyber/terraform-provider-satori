package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
	"sort"
	"strings"
)

// converting API objects back&forth to TF objects
// the function resolves the limitation that TypeMap can't have the TypeList elements
// !!! tf resources and api structures have to be similiar !!!
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
