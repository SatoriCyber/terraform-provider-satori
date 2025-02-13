package resources

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"strings"
)

// This function makes a consistent string that represents json.
// A map is created and when it is being serialized the order of the keys is
// defined by the json.Marshal function including the whitespaces, indentations etc.
// This causes these 2 json objects that are represented as strings to be equal i.e. normalized
// {"name": ron,       age:30} and {age:30,"name:ron}
func normalizeDataJSON(val interface{}) string {
	dataMap := map[string]interface{}{}

	// Ignoring errors since we know it is valid
	_ = json.Unmarshal([]byte(val.(string)), &dataMap)
	ret, _ := json.Marshal(dataMap)

	return string(ret)
}

func StringIsNotWhiteSpaceInArray(v interface{}, p cty.Path) diag.Diagnostics {
	value, ok := v.(string)
	var diags diag.Diagnostics

	if !ok {
		diagValue := diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Wrong value type",
			Detail:        fmt.Sprintf("expected type of %s to be string", p),
			AttributePath: p,
		}
		diags = append(diags, diagValue)
	} else {
		if strings.TrimSpace(value) == "" {
			errorMessage := fmt.Sprintf("value is expected to not be an empty string or whitespace.")

			attr, okIndexStepCasting := p[len(p)-1].(cty.IndexStep)
			if okIndexStepCasting && attr.Key.AsBigFloat() != nil {
				index := attr.Key.AsBigFloat().String()
				errorMessage = fmt.Sprintf("value at index %s is expected to not be an empty string or whitespace.", index)
			}

			diagValue := diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Empty string is not allowed",
				Detail:        errorMessage,
				AttributePath: p,
			}
			diags = append(diags, diagValue)
		}
	}
	return diags
}
