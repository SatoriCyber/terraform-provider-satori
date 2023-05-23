package resources

import (
	"encoding/json"
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
