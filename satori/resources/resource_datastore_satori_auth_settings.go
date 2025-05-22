package resources

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"log"
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
						}},
				},
				EnablePersonalAccessToken: &schema.Schema{
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Enables Satori Personal Access Token authentication for this data store. to be able using personal access token for authentication on this data store - data store temporary credentials must be enabled and personal access token feature should be enabled for the current account (see Account setting page in Satori platform).",
					Default:     false,
				},
			},
		},
	}
}
func GetSatoriAuthSettingsDatastoreOutput(d *schema.ResourceData, result *api.DataStoreOutput, err error) (map[string]interface{}, error) {
	var inInterface map[string]interface{}
	inJson, _ := json.Marshal(result.SatoriAuthSettings)
	err = json.Unmarshal(inJson, &inInterface)
	if err != nil {
		return nil, err
	}
	tfMap := biTfApiConverter(inInterface, false)
	if len(tfMap) == 0 { // empty result, skip it.
		return tfMap, nil
	}

	// If arrived here, it means that there is a configuration...
	// Check if the password has changed.
	// If the password has not changed, we need to set the password to the old/new value to state (simulating the backend response)
	// This is done for terraform update bypass.
	// The implementation is based on the fact that the password is stored in the terraform state.
	passwordResourcePath := "satori_auth_settings.0.credentials.0.password"
	if !d.HasChange(passwordResourcePath) { // no change
		log.Printf("The password hasn't change from state, overriding it with the old value")
		oldV, _ := d.GetChange(passwordResourcePath)
		credentialsMap := tfMap[Credentials].([]map[string]interface{})
		if len(credentialsMap) > 0 { // found credentials object (has to be defined)
			credentials := credentialsMap[0]      // it can be only 1 credentials object
			credentials[Password] = oldV.(string) // override the password with the old value
		}
	}

	return tfMap, nil
}

func SatoriAuthSettingsToResource(d *schema.ResourceData, in []interface{}) (*api.SatoriAuthSettings, error) {
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
