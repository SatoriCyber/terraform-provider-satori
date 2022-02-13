package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
	"sort"
	"strings"
)

// converting API generated basepolicy to terraform friendly map
func deepCopyMap(m map[string]interface{}, camelCase bool) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		vd, okVd := v.([]interface{})
		if (v) == nil && !okVd {
			cp[resNameTfConvert(k, camelCase)] = nil
		} else if ok {
			cp[resNameTfConvert(k, camelCase)] = []map[string]interface{}{deepCopyMap(vm, camelCase)}
		} else if okVd {
			var cd []map[string]interface{}
			for _, s := range vd {
				if s != nil {
					if currVal := deepCopyMap(s.(map[string]interface{}), camelCase); currVal != nil {
						cd = append(cd, currVal)
					}
				}
			}
			if !TreatAsMap[k] {
				cp[resNameTfConvert(k, camelCase)] = cd
			} else {
				var mpSa []interface{}
				for _, curRecord := range cd {
					for _, vaa := range curRecord {
						mpSa = append(mpSa, vaa)
					}
				}
				if len(cd) != 0 {
					cp[resNameTfConvert(k, camelCase)] = cd[0]
				} else {
					cp[resNameTfConvert(k, camelCase)] = []map[string]interface{}{}
				}
			}
		} else {
			cp[resNameTfConvert(k, camelCase)] = v
		}
	}
	return cp
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
