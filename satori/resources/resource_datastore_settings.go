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
		Description: "Settings for a MongoDB Data Store type",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				DeploymentType: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "MongoDB deployment type, for now supports only mongodb+srv and mongodb deployment",
				},
				AwsHostedZoneId: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "MongoDB AWS Hosted Zone ID, The Hosted AWS DNS Zone created for mapping MongoDB SRV records to Satori.",
				},
				AwsServerRoleArn: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "MongoDB AWS Service Role ARN, The IAM role ARN assumed by the DAC and used for updating records in the hosted DNS zone.",
				},
			},
		},
	}
}

func GetDatabricksSettingsDefinition() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Settings for a Databricks Data Store type",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				DatabricksAccountId: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Account ID",
				},
				DatabricksWarehouseId: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "SQL Warehouse ID",
				},
				Credentials: GetDatabricksCredentialsDefinition(),
			},
		},
	}
}

func GetDatabricksCredentialsDefinition() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "Credentials for Databricks Data Store type",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				Type: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Credentials type, user `AWS_SERVICE_PRINCIPAL_TOKEN` for AWS Service Principal Authentication",
				},
				DatabricksClientId: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Application (client) ID",
				},
				DatabricksClientSecret: {
					Type:        schema.TypeString,
					Sensitive:   true,
					Required:    true,
					Description: "Client secret value",
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

func GetSatoriDatastoreSettingsDatastoreOutput(result *api.DataStoreOutput, err error) (map[string]interface{}, error) {
	var inInterface map[string]interface{}
	inJson, _ := json.Marshal(result.DataStoreSettings)
	err = json.Unmarshal(inJson, &inInterface)
	if err != nil {
		return nil, err
	}
	tfMap := biTfApiConverter(inInterface, false)
	return tfMap, nil
}
