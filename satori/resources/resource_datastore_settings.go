package resources

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

func GetDataStoreSettingsDefinition() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Settings for a Data Store (may be unique per Data Store)",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				DeploymentType: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "MongoDB deployment type, for now supports only mongodb+srv and mongodb deployment",
				},
			},
		},
	}
}

func DataStoreSettingsToResource(in []interface{}) (*api.DataStoreSettings, error) {
	var dataStoreSettings api.DataStoreSettings
	mapDataStoreSettings := extractMapFromInterface(in)
	if mapDataStoreSettings != nil {
		tfMap := biTfApiConverter(mapDataStoreSettings, true)
		jsonOutput, _ := json.Marshal(tfMap)
		err := json.Unmarshal(jsonOutput, &dataStoreSettings)
		if err != nil {
			return nil, err
		}
	}
	return &dataStoreSettings, nil
}
