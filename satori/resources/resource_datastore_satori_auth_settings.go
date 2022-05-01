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
		Description: "Sets temporary credentials for admin to creeate temporary user datastore",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				Enabled: &schema.Schema{
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Enables to activate the creation of temporary credentials for this data store.",
					Default:     false,
				},
				Credentials: {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Credentials for Satori User Admin",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							Username: &schema.Schema{
								Type:        schema.TypeString,
								Required:    true,
								Description: "An admin username with rights to create a new user",
							},
							Password: &schema.Schema{
								Type:        schema.TypeString,
								Sensitive:   true,
								Required:    true,
								Description: "Password for the admin user",
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
