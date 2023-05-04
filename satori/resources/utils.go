package resources

import (
	"encoding/json"
	"os"
)

func normalizeDataJSON(val interface{}) string {
	dataMap := map[string]interface{}{}
	fileContent, readFileErr := os.ReadFile("/Users/rontzabary/workspace/local-provider-scripts/test.json")

	if readFileErr == nil {
		val = string(fileContent)
	}

	// Ignoring errors since we know it is valid
	_ = json.Unmarshal([]byte(val.(string)), &dataMap)
	ret, _ := json.Marshal(dataMap)

	return string(ret)
}
