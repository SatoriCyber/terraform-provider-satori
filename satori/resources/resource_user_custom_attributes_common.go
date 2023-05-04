package resources

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"reflect"
	"strings"
)

var (
	Attributes = "attributes"
	UserId     = "user_id"
)

func getUserCustomAttributesDefinitionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		UserId: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "User ID to manage Satori attributes for.",
		},
		Attributes: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "User's set of attributes in raw JSON object format. may include the following types: int, string, float, boolean, string[], number[], where number may be float/int  for example:  \"{\"company\": \"SatoriCyber\",\"age\": 30.5,\"cities\": [\"Washington\", \"Lisbon\"],\"kids_age\": [1, 3.14759, 7], \"isActive\": true}\" ",
			ValidateDiagFunc: func(i interface{}, k cty.Path) diag.Diagnostics {
				var unparsedJsonString map[string]interface{}

				v, ok := i.(string)

				if !ok {
					return diag.Errorf("Expected type %s to be string", k)
				}

				fileContent, readFileErr := os.ReadFile(v)

				if readFileErr == nil {
					v = string(fileContent)
				}

				// Validate that the user enters a map format and not just a json node e.g. string/number/boolean...
				err := json.Unmarshal([]byte(v), &unparsedJsonString)

				if err != nil {
					return diag.Errorf("Failed to parse rawJSON string at path %s, with err %s", k, err)
				}
				// TODO: change the validation method to the appropriate 1
				if !validateMapType(unparsedJsonString) {
					return diag.Errorf("raw JSON format is not valid OR not a json object, should be an attributes map, for example: \"{\"age\": 30, ...}\"")
				}

				if !validMapElementsAttributesType(unparsedJsonString) {
					return diag.Errorf("The raw JSON object contains an invalid value types.. valid types are { string, int, float, []string, []number } where number may be int|float")
				}

				return nil
			},
			StateFunc: normalizeDataJSON,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				if new == "" {
					return true
				}

				var oldCustomAttrs map[string]interface{}
				_ = json.Unmarshal([]byte(old), &oldCustomAttrs)

				var newCustomAttrs map[string]interface{}
				_ = json.Unmarshal([]byte(new), &newCustomAttrs)

				return reflect.DeepEqual(oldCustomAttrs, newCustomAttrs)
			},
		},
	}
}

func validateMapType(i interface{}) bool {
	typeOfObject := fmt.Sprintf("%T", i)
	if !strings.HasPrefix(typeOfObject, "map") {
		return false
	}
	return true
}

func validMapElementsAttributesType(m map[string]interface{}) bool {
	for _, value := range m {
		t := fmt.Sprintf("%T", value)
		if strings.HasPrefix(t, "[]") {
			if !validListTypeElement(value.([]interface{}), t) {
				return false
			}
			continue
		}
		if t != "bool" && t != "string" && !strings.HasPrefix(t, "int") && !strings.HasPrefix(t, "float") {
			return false
		}
	}

	return true
}

// This function validates that an element of type list will have a valid type
// Valid types are: []string, []int, []float, []float|int
func validListTypeElement(in []interface{}, t string) bool {

	if !strings.HasPrefix(t, "[]string") &&
		!strings.HasPrefix(t, "[]int") &&
		!strings.HasPrefix(t, "[]float") &&
		!strings.HasPrefix(t, "[]interface") {
		return false
	}

	if !strings.HasPrefix(t, "[]interface") {
		return false
	}

	allNumbers := true
	allStrings := true

	// Checking if all the values in the list is strings OR if all are numbers
	for _, val := range in {
		valType := fmt.Sprintf("%T", val)
		if !strings.HasPrefix(valType, "int") && !strings.HasPrefix(valType, "float") {
			allNumbers = false
		}
		if !strings.HasPrefix(valType, "string") {
			allStrings = false
		}
	}

	if !allStrings && !allNumbers {
		return false
	}
	return true
}
