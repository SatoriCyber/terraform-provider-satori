package resources

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

func GetSatoriAuthSettingsDefinitions() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Sets the authentication settings for the Data Store",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				Enabled: &schema.Schema{
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Enables Satori Data Store authentication.",
					Default:     false,
				},
				Credentials: {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Root user credentials",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							Username: &schema.Schema{
								Type:        schema.TypeString,
								Required:    true,
								Description: "Username of root user",
							},
							Password: &schema.Schema{
								Type:      schema.TypeString,
								Sensitive: true,
								Required:  true,
								Description: "Password of root user. This property is sensitive, and API does not return it in output. " +
									"In order to bypass terraform update, use lifecycle.ignore_changes, see example.",
							},
						}}},
			},
		},
	}
}
func GetSatoriAuthSettingsDatastoreOutput(result *api.DataStoreOutput, err error) (map[string]interface{}, error) {
	var inInterface map[string]interface{}
	inJson, _ := json.Marshal(result.SatoriAuthSettings)
	err = json.Unmarshal(inJson, &inInterface)
	if err != nil {
		return nil, err
	}
	tfMap := biTfApiConverter(inInterface, false)
	return tfMap, nil
}

func SatoriAuthSettingsToResource(in []interface{}) (*api.SatoriAuthSettings, error) {
	var satoriAuthSettings api.SatoriAuthSettings
	mapSatoriAuthSettings := extractMapFromInterface(in)
	if mapSatoriAuthSettings != nil {
		tfMap := biTfApiConverter(mapSatoriAuthSettings, true)
		jsonOutput, _ := json.Marshal(tfMap)
		err := json.Unmarshal(jsonOutput, &satoriAuthSettings)
		if err != nil {
			return nil, err
		}
	}
	return &satoriAuthSettings, nil
}
