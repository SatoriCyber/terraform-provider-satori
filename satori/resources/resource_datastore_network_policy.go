package resources

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

func GetNetworkPolicyDefinition() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
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
								Description: "Description",
							},
							IpRanges: GetIPRangesSchemaDefinitions(),
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
								Description: "Description",
							},
							IpRanges: GetIPRangesSchemaDefinitions(),
						},
					},
				},
			},
		},
	}
}

func GetIPRangesSchemaDefinitions() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "IP Ranges",
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
	if mapNetworkPolicy == nil {
		return nil, errors.New("networkPolicy is incorrect/missing")
	}
	tfMap := biTfApiConverter(mapNetworkPolicy, true)
	jsonOutput, _ := json.Marshal(tfMap)
	err := json.Unmarshal(jsonOutput, &networkPolicy)
	if err != nil {
		return nil, err
	}
	return &networkPolicy, nil
}
