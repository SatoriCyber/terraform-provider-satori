package resources

import (
	"encoding/json"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

func GetBaseLinePolicyDefinition() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		MaxItems:    1,
		Description: "Baseline security policy.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				Type: &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Description: "DataStore basepolicy .",
					Default:     "BASELINE_POLICY",
				},
				UnassociatedQueriesCategory: {
					Type:        schema.TypeList,
					Required:    true,
					MaxItems:    1,
					Description: "UnassociatedQueriesCategory",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							QueryAction: &schema.Schema{
								Type:        schema.TypeString, //Default:     "PASS",
								Optional:    true,
								Description: "Default policy action for querying locations that are not associated with a dataset, modes supported:  PASS┃REDACT┃BLOCK.",
							}}}},

				UnsupportedQueriesCategory: {
					Type:        schema.TypeList,
					Required:    true,
					MaxItems:    1,
					Description: "UnsupportedQueriesCategory",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							QueryAction: &schema.Schema{
								Type:        schema.TypeString, //Default:     "PASS",
								Required:    true,
								Description: "Default policy action for unsupported queries and objects, modes supported:  PASS┃REDACT┃BLOCK ",
							}}}},

				Exclusions: {
					Type:        schema.TypeList,
					Required:    true,
					MaxItems:    1,
					Description: "Exempt users and patterns from baseline security policy",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							ExcludedIdentities: &schema.Schema{
								Type:        schema.TypeList,
								Optional:    true,
								Default:     nil,
								Description: "Exempt Users from the Baseline Security Policy\n",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										IdentityType: &schema.Schema{
											Type:        schema.TypeString,
											Optional:    true,
											Description: "USER type is supported",
										}, Identity: &schema.Schema{
											Type:        schema.TypeString,
											Optional:    true,
											Description: "Username",
										},
									}},
							},
							ExcludedQueryPatterns: &schema.Schema{
								Type:        schema.TypeList,
								Optional:    true,
								Default:     nil,
								Description: "Exempt Queries from the Baseline Security Policy",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										Pattern: &schema.Schema{
											Type:        schema.TypeString,
											Optional:    true,
											Description: "Query pattern",
										},
									}},
							},
						},
					}},
			},
		},
	}
}

func GetBaseLinePolicyDatastoreOutput(result *api.DataStoreOutput, err error) (map[string]interface{}, error) {
	var inInterface map[string]interface{}
	inJson, _ := json.Marshal(result.BaselineSecurityPolicy)
	err = json.Unmarshal(inJson, &inInterface)
	if err != nil {
		return nil, err
	}
	tfMap := biTfApiConverter(inInterface, false)
	return tfMap, nil
}

func BaselineSecurityPolicyToResource(in []interface{}) (*api.BaselineSecurityPolicy, error) {
	var baselineSecurityPolicy api.BaselineSecurityPolicy
	mapBaselineSecurityPolicy := extractMapFromInterface(in)
	if mapBaselineSecurityPolicy == nil {
		return nil, errors.New("datastore is incorrect/missing")
	}
	tfMap := biTfApiConverter(mapBaselineSecurityPolicy, true)
	jsonOutput, _ := json.Marshal(tfMap)
	err := json.Unmarshal(jsonOutput, &baselineSecurityPolicy)
	if err != nil {
		return nil, err
	}
	return &baselineSecurityPolicy, nil
}
