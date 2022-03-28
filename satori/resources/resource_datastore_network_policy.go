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
				Name: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Policy name - ACL",
				},
				AllowedRules: {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Allowed Ip Rules",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							IPRanges: GetIPRangesSchemaDefinitions(),
						},
					},
				},
				BlockedRules: {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Blocked Ips Rules",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							IPRanges: GetIPRangesSchemaDefinitions(),
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
				Note: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Description",
				},
				IPRange: {
					Type:        schema.TypeList,
					Required:    true,
					MinItems:    1,
					Description: "IP Range",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func GetNetworkPolicyDatastoreOutput(result *api.NetworkPolicy, err error) (map[string]interface{}, error) {
	var inInterface map[string]interface{}
	inJson, _ := json.Marshal(result)
	err = json.Unmarshal(inJson, &inInterface)
	if err != nil {
		return nil, err
	}
	tfMap := biTfApiConverter(inInterface, false)

	//if allowedRules, ok := tfMap[AllowedRules].([]map[string]interface{}); ok {
	//	for _, allowedRule := range allowedRules {
	//    var note = allowedRule[Note].(string)
	//		if ipRanges, ok := allowedRule[IPRanges].([]map[string]interface{}); ok {
	//			var newArr []string
	//			for _, ipRangesRow := range ipRanges {
	//        ipRangesRow[Note] = note
	//				if ipRange, ok := ipRangesRow[IPRange].(string); ok {
	//					newArr = append(newArr, ipRange)
	//				}
	//			}
	//			allowedRule[IPRanges] = newArr
	//			newArr = nil
	//		}
	//	}
	//}
	//
	//if blockedRules, ok := tfMap[BlockedRules].([]map[string]interface{}); ok {
	// for _, blockedRule := range blockedRules {
	//   if ipRanges, ok := blockedRule[IPRanges].([]map[string]interface{}); ok {
	//     var newArr []string
	//     for _, ipRangesRow := range ipRanges {
	//       if ipRange, ok := ipRangesRow[IPRange].(string); ok {
	//         newArr = append(newArr, ipRange)
	//       }
	//     }
	//     blockedRule[IPRanges] = newArr
	//     newArr = nil
	//   }
	// }
	//}

	return tfMap, nil
}

func NetworkPolicyToResource(in []interface{}) (*api.NetworkPolicy, error) {
	var out api.NetworkPolicy
	var values = extractMapFromInterface(in)

	if name, ok := values[Name].(string); ok {
		out.Name = name
	}
	if allowedRules, ok := values[AllowedRules]; ok {
		var allowedRules, err1 = NetworkPolicyRulesToResource(allowedRules.([]interface{}))
		if err1 != nil {
			return &out, err1
		}
		out.AllowedRules = allowedRules
	}
	if blockedRules, ok := values[BlockedRules]; ok {
		var blockedRules, err2 = NetworkPolicyRulesToResource(blockedRules.([]interface{}))
		if err2 != nil {
			return &out, err2
		}
		out.BlockedRules = blockedRules
	}

	return &out, nil
}

func NetworkPolicyRulesToResource(values []interface{}) ([]api.NetworkPolicyRules, error) {
	var out []api.NetworkPolicyRules
	for i := range values {
		if value, ok := values[i].(map[string]interface{}); ok {
			if ipRangesSource, ok := value[IPRanges]; ok {
				var ipRangesArr = ipRangesSource.([]interface{})
				var ipRangesArrayOut []api.IPRanges = nil
				for _, ipRanges := range ipRangesArr {

					// ips
					if ipRangeSource, ok := ipRanges.(map[string]interface{})[IPRange]; ok {
						for _, ip := range ipRangeSource.([]interface{}) {
							var ipRangesOut api.IPRanges
							ipRangesOut.IPRange = ip.(string)
							ipRangesArrayOut = append(ipRangesArrayOut, ipRangesOut)
						}
					}

					// build the NetworkPolicyRules record
					var o api.NetworkPolicyRules
					o.Note = ipRanges.(map[string]interface{})[Note].(string)
					o.IPRanges = ipRangesArrayOut
					ipRangesArrayOut = nil
					out = append(out, o)
				}
			}
		}
	}
	return out, nil
}
