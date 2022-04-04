package resources

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

func GetNetworkPolicyDefinition() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "a Network Policy for a Data Store",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				AllowedRules: {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Allowed Ip Rules",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							Note: {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "custom description for allowed IP ranges",
							},
							IpRanges: GetIpRangesSchemaDefinitions(),
						},
					},
				},
				BlockedRules: {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Blocked Ips Rules",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							Note: {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "custom description for blocked IP ranges",
							},
							IpRanges: GetIpRangesSchemaDefinitions(),
						},
					},
				},
			},
		},
	}
}

func GetIpRangesSchemaDefinitions() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "enable access control from specified IP ranges",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				IpRange: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "IP Range",
				},
			},
		},
	}
}

func GetNetworkPolicyDatastoreOutput(result *api.DataStoreOutput, err error) (map[string]interface{}, error) {
	var inInterface map[string]interface{}
	inJson, _ := json.Marshal(result.NetworkPolicy)
	err = json.Unmarshal(inJson, &inInterface)
	if err != nil {
		return nil, err
	}
	tfMap := biTfApiConverter(inInterface, false)
	return tfMap, nil
}

func NetworkPolicyToResource(in []interface{}) (*api.NetworkPolicy, error) {
	var networkPolicy api.NetworkPolicy
	mapNetworkPolicy := extractMapFromInterface(in)
	if mapNetworkPolicy != nil {
		tfMap := biTfApiConverter(mapNetworkPolicy, true)
		jsonOutput, _ := json.Marshal(tfMap)
		err := json.Unmarshal(jsonOutput, &networkPolicy)
		if err != nil {
			return nil, err
		}
	}
	return &networkPolicy, nil
}
