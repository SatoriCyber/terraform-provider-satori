package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"log"
)

var (
	Active                             = "active"
	SecurityPolicyName                 = "name"
	SecurityPolicyProfile              = "profile"
	MaskingProfile                     = "masking"
	MaskingRule                        = "rule"
	RuleId                             = "id"
	RuleDescription                    = "description"
	MaskingRuleAction                  = "action"
	MaskingRuleActionType              = "type"
	MaskingRuleActionProfileId         = "masking_profile_id"
	RuleCriteria                       = "criteria"
	CriteriaCondition                  = "condition"
	CriteriaIdentity                   = "identity"
	MaskingRuleConditionMasking        = "conditional_masking"
	MaskingRuleWhereCondition          = "where_condition"
	MaskingRuleActionDefaultActionType = "APPLY_MASKING_PROFILE"
	RowLevelSecurity                   = "row_level_security"
	RLSActive                          = "active"
	RLSRule                            = "rule"
	RLSRuleFilter                      = "filter"
	RLSRuleFilterDatastoreId           = "datastore_id"
	RLSRuleFilterLocationPrefix        = "location" // deprecated , This is the deprecated field and should be removed
	RLSRuleFilterAdvanced              = "advanced"
	RLSRuleFilterLogicYaml             = "logic_yaml"
	RLSMapping                         = "mapping"
	FilterName                         = "name"
	RLSMappingFilter                   = "filter"
	RLSMappingValue                    = "value"
	RLSMappingValues                   = "values"
	RLSMappingDefaultValues            = "defaults"
	RLSMappingValuesType               = "type"
)

func ResourceSecurityPolicy() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceSecurityPolicyCreate,
		ReadContext:   resourceSecurityPolicyRead,
		UpdateContext: resourceSecurityPolicyUpdate,
		DeleteContext: resourceSecurityPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Security Policy.",
		Schema: map[string]*schema.Schema{
			SecurityPolicyName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Security policy name.",
			},
			SecurityPolicyProfile: {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Security policy profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						MaskingProfile:   resourceMaskingProfile(),
						RowLevelSecurity: resourceRowLevelSecurity(),
					},
				},
			},
		},
	}
}

func resourceRowLevelSecurity() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Row level security profile",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				RLSActive: {
					Type:        schema.TypeBool,
					Required:    true,
					Description: "Row level security activation.",
				},
				RLSRule:    resourceRowLevelSecurityRule(),
				RLSMapping: resourceRowLevelSecurityMappingFilter(),
			},
		},
	}
}

func resourceRowLevelSecurityMappingFilter() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Row Level Security Mapping.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				FilterName: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Filter name, has to be unique in this policy.",
				},
				RLSMappingFilter: {
					Type:        schema.TypeList,
					Required:    true,
					MaxItems:    1,
					Description: "Filter definition.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							RuleCriteria: {
								Type:        schema.TypeList,
								Required:    true,
								MinItems:    1,
								MaxItems:    1,
								Description: "Filter criteria.",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										CriteriaCondition: {
											Type:        schema.TypeString,
											Required:    true,
											Description: "Identity condition, for example IS_NOT, IS, etc.",
										},
										CriteriaIdentity: resourceDataAccessIdentity(true),
									},
								},
							},
							RLSMappingValues: {
								Type:        schema.TypeList,
								Required:    true,
								MaxItems:    1,
								Description: "A list of values to be applied in this filter. Values are dependent on their type and has to be homogeneous",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										RLSMappingValue: {
											Type:        schema.TypeList,
											Required:    true,
											MinItems:    1,
											Description: "List of values, when ANY_VALUE or ALL_OTHER_VALUES are defined, the list has to be empty",
											Elem: &schema.Schema{
												Type: schema.TypeString,
											},
										},
										RLSMappingValuesType: {
											Type:        schema.TypeString,
											Required:    true,
											Description: "Values type. Allowed options: STRING, NUMERIC, CEL, SQL, ANY_VALUE, ALL_OTHER_VALUES.",
											ValidateFunc: func(v interface{}, key string) (warns []string, errs []error) {
												value := v.(string)
												if value != "STRING" && value != "NUMERIC" && value != "CEL" && value != "SQL" && value != "ANY_VALUE" && value != "ALL_OTHER_VALUES" {
													errs = append(errs, fmt.Errorf("%q must be one of 'STRING, NUMERIC, CEL, SQL, ANY_VALUE or ALL_OTHER_VALUES' but got: %q", key, value))
												}
												return
											},
										},
									},
								},
							},
						},
					},
				},
				RLSMappingDefaultValues: {
					Type:        schema.TypeList,
					Required:    true,
					MaxItems:    1,
					MinItems:    1,
					Description: "A list of default values to be applied in this filter if there was no match. Values are dependent on their type and has to be homogeneous",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							RLSMappingValue: {
								Type:        schema.TypeList,
								Required:    true,
								Description: "List of values, when NO_VALUE or ALL_OTHER_VALUES are defined, the list has to be empty",
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
							RLSMappingValuesType: {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Default values type. Allowed options: STRING, NUMERIC, CEL, SQL, NO_VALUE, ALL_OTHER_VALUES.",
								ValidateFunc: func(v interface{}, key string) (warns []string, errs []error) {
									value := v.(string)
									if value != "STRING" && value != "NUMERIC" && value != "CEL" && value != "SQL" && value != "NO_VALUE" && value != "ALL_OTHER_VALUES" {
										errs = append(errs, fmt.Errorf("%q must be one of 'STRING, NUMERIC, CEL, SQL, NO_VALUE or ALL_OTHER_VALUES' but got: %q", key, value))
									}
									return
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceRowLevelSecurityRule() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Row Level Security Rule definition.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				RuleId: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Rule id, has to be unique.",
				},
				RuleDescription: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Rule description.",
				},
				Active: {
					Type:        schema.TypeBool,
					Required:    true,
					Description: "Is active rule.",
				},
				RLSRuleFilter: {
					Type:        schema.TypeList,
					Required:    true,
					MaxItems:    1,
					Description: "Rule filter.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							RLSRuleFilterDatastoreId: {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Datastore ID.",
							},
							RLSRuleFilterLocationPrefix: {
								Type:        schema.TypeList,
								Optional:    true,
								Deprecated:  "The 'location' field has been deprecated. Please use the 'location_path', `location_parts` or `location_parts_full` fields instead.",
								Description: "Location to be included in the rule.",
								Elem:        getLocationResource(),
							},
							LocationPath: {
								Type:             schema.TypeString,
								Optional:         true,
								Description:      "The short presentation of the location path in the data store. Includes `.` separated string when part types are defined with default definitions. For example 'a.b.c' in Snowflake data store will path to table 'a' under schema 'b' under database 'a'.  Conflicts with 'location', 'location_parts', and 'location_parts_full' fields.",
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
							},
							LocationParts: {
								Type:        schema.TypeList,
								Optional:    true,
								MinItems:    1,
								Description: "The part separated location path in the data store. Includes an array of path parts when part types are defined with default definitions. For example ['a', 'b', 'c'] in Snowflake data store will path to table 'a' under schema 'b' under database 'a'. Conflicts with 'location', 'location_path', and 'location_parts_full' fields",
								Elem: &schema.Schema{
									Type:             schema.TypeString,
									ValidateDiagFunc: StringIsNotWhiteSpaceInArray,
								},
							},
							LocationPartsFull: {
								Type:        schema.TypeList,
								Optional:    true,
								Description: "The full location path definition in the data store. Includes an array of objects with path name and path type. Can be used when the path type should be defined explicitly and not as defined by default. For example [{name= 'a', type= 'DATABASE'},{name= 'b', type= 'SCHEMA'},{name= 'view.c', type= 'VIEW'}]. Conflicts with 'location', 'location_path', and 'location_parts' fields.",
								MinItems:    1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"name": {
											Type:             schema.TypeString,
											Required:         true,
											Description:      "The name of the location part.",
											ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
										},
										"type": {
											Type:             schema.TypeString,
											Required:         true,
											Description:      "The type of the location part. Optional values: TABLE, COLUMN, SEMANTIC_MODEL, REPORT, DASHBOARD, DATABASE, SCHEMA, JSON_PATH, WAREHOUSE, ENDPOINT, TYPE, FIELD, EXTERNAL_LOCATION, CATALOG, BUCKET, OBJECT, COLLECTION, VIEW, etc",
											ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
										},
									},
								},
							},
							RLSRuleFilterAdvanced: {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     true,
								Description: "Describes if logic yaml contains complex configuration.",
							},
							RLSRuleFilterLogicYaml: {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Conditional rule, for more info see https://satoricyber.com/docs/security-policies/#setting-up-data-filtering.",
							},
						},
					},
				},
			},
		},
	}
}

func resourceMaskingProfile() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Masking profile.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				Active: {
					Type:        schema.TypeBool,
					Required:    true,
					Description: "Is active.",
				},
				MaskingRule: {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Masking Rule.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							RuleId: {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Rule id, has to be unique.",
							},
							RuleDescription: {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Rule description.",
							},
							Active: {
								Type:        schema.TypeBool,
								Required:    true,
								Description: "Is active rule.",
							},
							MaskingRuleAction: {
								Type:        schema.TypeList,
								Required:    true,
								MaxItems:    1,
								Description: "Rule action.",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										MaskingRuleActionType: {
											Type:        schema.TypeString,
											Optional:    true,
											Default:     MaskingRuleActionDefaultActionType,
											Description: "Rule type.",
										},
										MaskingRuleActionProfileId: {
											Type:        schema.TypeString,
											Required:    true,
											Description: "The reference id to be applied as masking profile.",
										},
									},
								},
							},
							RuleCriteria: {
								Type:        schema.TypeList,
								Required:    true,
								MinItems:    1,
								MaxItems:    1,
								Description: "Masking criteria.",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										CriteriaCondition: {
											Type:        schema.TypeString,
											Required:    true,
											Description: "Identity condition, for example IS_NOT, IS, etc.",
										},
										CriteriaIdentity: resourceDataAccessIdentity(true),
									},
								},
							},
							MaskingRuleConditionMasking: {
								Type:        schema.TypeList,
								MinItems:    1,
								MaxItems:    1,
								Optional:    true,
								Description: "Conditional masking. Only supported in the Databricks and Snowflake Native Integrations.",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										MaskingRuleWhereCondition: {
											Type:        schema.TypeString,
											Required:    true,
											Description: "Where condition.",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// //////////////////////////////////
// Resource to schema mappers
// //////////////////////////////////
func resourceToSecurityProfiles(d *schema.ResourceData) (*api.SecurityProfiles, error) {
	if _, ok := d.GetOk(SecurityPolicyProfile); !ok {
		return nil, nil
	}

	out := api.SecurityProfiles{}

	if _, ok := d.GetOk("profile.0.masking"); ok {
		if m, ok := d.GetOk("profile.0.masking.0"); ok {
			out.Masking = &api.MaskingSecurityProfile{}
			resourceToMasking(d, m, &out)
		}
	}

	if _, ok := d.GetOk("profile.0.row_level_security"); ok {
		if m, ok := d.GetOk("profile.0.row_level_security.0"); ok {
			out.RowLevelSecurity = &api.RowLevelSecurityProfile{}
			err := resourceToRowLevelSecurityProfile(d, m, &out)
			if err != nil {
				return nil, err
			}
		}
	}

	return &out, nil
}

// Masking
func resourceToMasking(d *schema.ResourceData, m interface{}, out *api.SecurityProfiles) {
	masking := m.(map[string]interface{})

	isActive := masking[Active].(bool)

	log.Printf("Masking is active: %t", isActive)
	out.Masking.Active = isActive

	if v, ok := d.GetOk("profile.0.masking.0.rule"); ok {
		rules := make([]api.MaskingRule, len(v.([]interface{})))
		for i, raw := range v.([]interface{}) {
			resourceToMaskingRule(raw, &rules, i)
		}
		out.Masking.Rules = rules
	}
}

func resourceToMaskingRule(raw interface{}, rules *[]api.MaskingRule, i int) {
	inElement := raw.(map[string]interface{})
	outElement := api.MaskingRule{}
	outElement.Id = inElement[RuleId].(string)
	outElement.Description = inElement[RuleDescription].(string)
	outElement.Active = inElement[Active].(bool)

	actionList := inElement[MaskingRuleAction].([]interface{})
	action := actionList[0].(map[string]interface{})
	maskingProfileId := action[MaskingRuleActionProfileId].(string)

	// Masking action
	outElement.MaskingAction.MaskingProfileId = maskingProfileId
	outElement.MaskingAction.Type = MaskingRuleActionDefaultActionType

	// Masking criteria
	outElement.DataFilterCriteria.Identity = *resourceToCriteriaRule(inElement, &outElement.DataFilterCriteria.Condition)

	// Conditional masking
	conditionalMaskingList := inElement[MaskingRuleConditionMasking].([]interface{})
	if len(conditionalMaskingList) > 0 {
		conditionalMasking := conditionalMaskingList[0].(map[string]interface{})
		outElement.ConditionalMasking.WhereCondition = conditionalMasking[MaskingRuleWhereCondition].(string)
	}

	(*rules)[i] = outElement
}

func resourceToCriteriaRule(inElement map[string]interface{}, condition *string) *api.DataAccessIdentity {
	criteriaList := inElement[RuleCriteria].([]interface{})
	criteria := criteriaList[0].(map[string]interface{})
	*condition = criteria[CriteriaCondition].(string)

	identityList := criteria[CriteriaIdentity].([]interface{})

	identity := identityList[0].(map[string]interface{})
	identityOut := resourceToIdentity(identity)

	return identityOut
}

// Row level security
func resourceToRowLevelSecurityProfile(d *schema.ResourceData, m interface{}, out *api.SecurityProfiles) error {
	rls := m.(map[string]interface{})

	isActive := rls[Active].(bool)
	out.RowLevelSecurity.Active = isActive

	if v, ok := d.GetOk("profile.0.row_level_security.0.rule"); ok {
		rules := make([]api.RowLevelSecurityRule, len(v.([]interface{})))
		for i, raw := range v.([]interface{}) {
			err := resourceToRowLevelSecurityRule(raw, &rules, i)
			if err != nil {
				return err
			}
		}
		out.RowLevelSecurity.Rules = rules
	}
	if v, ok := d.GetOk("profile.0.row_level_security.0.mapping"); ok {
		mapping := make([]api.RowLevelSecurityFilter, len(v.([]interface{})))
		for i, raw := range v.([]interface{}) {
			resourceToRowLevelSecurityFilter(raw, &mapping, i)
		}
		out.RowLevelSecurity.Maps = mapping
	}
	return nil
}

func resourceToRowLevelSecurityRule(raw interface{}, rules *[]api.RowLevelSecurityRule, i int) error {
	inElement := raw.(map[string]interface{})
	outElement := api.RowLevelSecurityRule{}

	outElement.Id = inElement[RuleId].(string)
	outElement.Description = inElement[RuleDescription].(string)
	outElement.Active = inElement[Active].(bool)

	filterList := inElement[RLSRuleFilter].([]interface{})
	filter := filterList[0].(map[string]interface{})

	datastoreId := filter[RLSRuleFilterDatastoreId].(string)
	logicYaml := filter[RLSRuleFilterLogicYaml].(string)
	advanced := filter[RLSRuleFilterAdvanced].(bool)

	// Masking action
	outElement.RuleFilter.DataStoreId = datastoreId
	outElement.RuleFilter.LogicYaml = logicYaml
	outElement.RuleFilter.Advanced = advanced

	//var location api.DataSetGenericLocation

	err := checkThatOnlyOneLocationFormatExists(filter, Location, LocationPath, LocationParts, LocationPartsFull, false)
	if err != nil {
		return err
	}

	filterPrint, _ := json.Marshal(filter)
	log.Printf("resourceToRowLevelSecurityRule: The filter presentation `%s`", filterPrint)

	if len(filter[Location].([]interface{})) > 0 { // deprecated field
		inLocations := filter[Location].([]interface{})
		if len(inLocations) > 0 {
			var location api.DataSetGenericLocation
			err := resourceToLocation(&location, inLocations, true)
			if err != nil {
				return err
			}
			outElement.RuleFilter.LocationPrefix = &location
		}
	} else if len(filter[LocationParts].([]interface{})) > 0 { // new field for LocationParts
		inLocations := filter[LocationParts].([]interface{})
		if len(inLocations) > 0 {
			var location []api.LocationPath
			err := resourcePartsToLocationPath(&location, inLocations, false)
			if err != nil {
				return err
			}
			outElement.RuleFilter.LocationPath = location
		}
	} else if len(filter[LocationPartsFull].([]interface{})) > 0 { // new field for LocationPartsFull
		inLocations := filter[LocationPartsFull].([]interface{})
		if len(inLocations) > 0 {
			var location []api.LocationPath
			err := resourcePartsFullToLocationPath(&location, inLocations, false)
			if err != nil {
				return err
			}
			outElement.RuleFilter.LocationPath = location
		}
	} else if filter[LocationPath] != nil && len(filter[LocationPath].(string)) > 0 { // new string field, ignore if empty
		// terraform value will be always not nil, so we need to check the length and consider it as does not exist if empty
		inLocationStr := filter[LocationPath].(string)

		var location []api.LocationPath
		log.Printf("resourceToRowLevelSecurityRule: found %s location path with length %d", inLocationStr, len(inLocationStr))
		err := resourceStrToLocationPath(&location, inLocationStr, false)
		if err != nil {
			return err
		}
		outElement.RuleFilter.LocationPath = location
	} else {
		return fmt.Errorf("resourceToRowLevelSecurityRule: location type is not defined")
	}

	outElementPrint, _ := json.Marshal(outElement)
	log.Printf("resourceToRowLevelSecurityRule: The outElement presentation `%s`", outElementPrint)

	(*rules)[i] = outElement

	return nil
}

func resourceToRowLevelSecurityFilter(raw interface{}, rules *[]api.RowLevelSecurityFilter, i int) {
	inElement := raw.(map[string]interface{})
	outElement := api.RowLevelSecurityFilter{}

	outElement.Name = inElement[FilterName].(string)

	// filter
	filterList := inElement[RLSMappingFilter].([]interface{})

	filters := make([]api.RowLevelSecurityMapDataFilter, len(filterList))
	for i, filter := range filterList {
		inFilter := filter.(map[string]interface{})
		var outFilter api.RowLevelSecurityMapDataFilter

		// Filter criteria
		outFilter.Criteria.Identity = *resourceToCriteriaRule(inFilter, &outFilter.Criteria.Condition)

		// Filter values
		valuesList := inFilter[RLSMappingValues].([]interface{})
		values := valuesList[0].(map[string]interface{})
		outFilter.Values.Type = values[RLSMappingValuesType].(string)

		valueIn := values[RLSMappingValue].([]interface{})
		if len(valueIn) > 0 {
			strValuesArray := make([]string, len(valueIn))
			for v, strValue := range valueIn {
				strValuesArray[v] = strValue.(string)
			}
			outFilter.Values.Values = &strValuesArray
		} else {
			outFilter.Values.Values = nil
		}

		filters[i] = outFilter
	}
	outElement.Filters = filters

	defaultsIn := inElement[RLSMappingDefaultValues].([]interface{})
	defaults := defaultsIn[0].(map[string]interface{})

	outElement.Defaults.Type = defaults[RLSMappingValuesType].(string)

	defaultValues := defaults[RLSMappingValue].([]interface{})
	if defaultValues != nil && len(defaultValues) > 0 {
		strDefaultValues := make([]string, len(defaultValues))
		for v, strDefaultValue := range defaultValues {
			strDefaultValues[v] = strDefaultValue.(string)
		}
		outElement.Defaults.Values = &strDefaultValues
	} else {
		outElement.Defaults.Values = nil
	}

	(*rules)[i] = outElement
}

func resourceToSecurityPolicy(d *schema.ResourceData) (*api.SecurityPolicy, error) {
	out := api.SecurityPolicy{}
	out.Name = d.Get(SecurityPolicyName).(string)

	securityPolicy, err := resourceToSecurityProfiles(d)
	if err != nil {
		return nil, err
	}
	out.SecurityProfiles = securityPolicy
	return &out, nil
}

func resourceSecurityPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	input, err := resourceToSecurityPolicy(d)
	if err != nil {
		log.Printf("Recieved an error in security policy mapping: %s", err)
		return diag.FromErr(err)
	}

	result, err := c.CreateSecurityPolicy(input)
	if err != nil {
		log.Printf("Recieved an error in security policy create: %s", err)
		return diag.FromErr(err)
	} else {
		d.SetId(result.Id)
	}

	return diags
}

// //////////////////////////////////
// Schema to resource mappers
// //////////////////////////////////
func securityProfilesToResource(profiles *api.SecurityProfiles, d *schema.ResourceData) interface{} {

	if profiles == nil {
		return nil
	}

	out := make([]map[string]interface{}, 1)

	out[0] = make(map[string]interface{})
	out[0][MaskingProfile] = maskingToResource(profiles.Masking)
	out[0][RowLevelSecurity] = rowLevelSecurityToResource(profiles.RowLevelSecurity, d)

	return out
}

func rowLevelSecurityToResource(security *api.RowLevelSecurityProfile, d *schema.ResourceData) interface{} {
	if security == nil {
		return nil
	}
	out := make([]map[string]interface{}, 1)
	out[0] = make(map[string]interface{})

	out[0][Active] = security.Active

	rules := make([]map[string]interface{}, len(security.Rules))
	mapping := make([]map[string]interface{}, len(security.Maps))

	for i, v := range security.Rules {
		rules[i] = make(map[string]interface{})
		rules[i][RuleId] = v.Id
		rules[i][Active] = v.Active
		rules[i][RuleDescription] = v.Description

		ruleFilter := make([]map[string]interface{}, 1)
		ruleFilter[0] = make(map[string]interface{})
		ruleFilter[0][RLSRuleFilterLogicYaml] = v.RuleFilter.LogicYaml
		ruleFilter[0][RLSRuleFilterDatastoreId] = v.RuleFilter.DataStoreId
		ruleFilter[0][RLSRuleFilterAdvanced] = v.RuleFilter.Advanced

		prefixFieldName := "profile.0.row_level_security.0.rule"

		if v.RuleFilter.LocationPrefix != nil {
			// Checks if the state already contains the deprecated field, if so, convert the output to the deprecated format,
			// otherwise convert to the new format
			log.Printf("rowLevelSecurityToResource: found old `location` deprecated format at, %s.%d.%s.%d.%s", prefixFieldName, i, RLSRuleFilter, 0, RLSRuleFilterLocationPrefix)
			ruleFilter[0][Location] = []map[string]interface{}{locationToResource(v.RuleFilter.LocationPrefix)}
		} else if v.RuleFilter.LocationPath != nil { // new field format
			if configuredLocationPath, ok := d.GetOk(fmt.Sprintf("%s.%d.%s.%d.%s", prefixFieldName, i, RLSRuleFilter, 0, LocationPath)); ok { // new LocationPath was configured
				log.Printf("rowLevelSecurityToResource: new format for %s was found, value: %s", LocationPath, configuredLocationPath)
				ruleFilter[0][LocationPath] = locationPathToLocationPathResource(v.RuleFilter.LocationPath)
			} else if configuredLocationPath, ok := d.GetOk(fmt.Sprintf("%s.%d.%s.%d.%s", prefixFieldName, i, RLSRuleFilter, 0, LocationParts)); ok { // new LocationPath was configured
				log.Printf("rowLevelSecurityToResource: new format for %s was found, value: %s", LocationParts, configuredLocationPath)
				ruleFilter[0][LocationParts] = locationPathToLocationPartsResource(v.RuleFilter.LocationPath)
			} else if configuredLocationPath, ok := d.GetOk(fmt.Sprintf("%s.%d.%s.%d.%s", prefixFieldName, i, RLSRuleFilter, 0, LocationPartsFull)); ok { // new LocationPath was configured
				log.Printf("rowLevelSecurityToResource: new format for %s was found, value: %s", LocationPartsFull, configuredLocationPath)
				ruleFilter[0][LocationPartsFull] = locationPathToLocationPartsFullResource(v.RuleFilter.LocationPath)
			} else {
				log.Printf("got an unknown format for locationPath")
			}
		}

		rules[i][RLSRuleFilter] = ruleFilter

	}

	for i, v := range security.Maps {
		mapping[i] = make(map[string]interface{})
		mapping[i][FilterName] = v.Name

		filters := make([]map[string]interface{}, len(v.Filters))
		for j, f := range v.Filters {
			filters[j] = make(map[string]interface{})

			criteria := make([]map[string]interface{}, 1)
			criteria[0] = make(map[string]interface{})
			criteria[0][CriteriaCondition] = f.Criteria.Condition

			identity := make([]map[string]interface{}, 1)
			identity[0] = *dataAccessIdentityToResource(&f.Criteria.Identity)
			criteria[0][CriteriaIdentity] = identity
			filters[j][RuleCriteria] = criteria

			values := make([]map[string]interface{}, 1)
			values[0] = make(map[string]interface{})
			values[0][RLSMappingValuesType] = f.Values.Type

			if f.Values.Values != nil && len(*f.Values.Values) > 0 {
				value := make([]string, len(*f.Values.Values))
				for k, strValue := range *f.Values.Values {
					value[k] = strValue
				}
				values[0][RLSMappingValue] = value
				filters[j][RLSMappingValues] = values
			}
		}
		mapping[i][RLSMappingFilter] = filters

		defaults := make([]map[string]interface{}, 1)
		defaults[0] = make(map[string]interface{})

		defaults[0][RLSMappingValuesType] = v.Defaults.Type

		if v.Defaults.Values != nil && len(*v.Defaults.Values) > 0 {
			defaultsValues := make([]string, len(*v.Defaults.Values))
			for k, strValue := range *v.Defaults.Values {
				defaultsValues[k] = strValue
			}
			defaults[0][RLSMappingValue] = defaultsValues
		} else {
			defaults[0][RLSMappingValue] = nil
		}

		mapping[i][RLSMappingDefaultValues] = defaults
	}

	out[0][RLSRule] = rules
	out[0][RLSMapping] = mapping
	return out
}

func maskingToResource(masking *api.MaskingSecurityProfile) interface{} {
	if masking == nil {
		return nil
	}
	out := make([]map[string]interface{}, 1)

	out[0] = make(map[string]interface{})
	out[0][Active] = masking.Active
	rules := make([]map[string]interface{}, len(masking.Rules))

	for i, v := range masking.Rules {
		rules[i] = make(map[string]interface{})
		rules[i][Active] = v.Active
		rules[i][RuleDescription] = v.Description
		rules[i][RuleId] = v.Id

		action := make([]map[string]interface{}, 1)
		action[0] = make(map[string]interface{})
		action[0][MaskingRuleActionType] = v.MaskingAction.Type
		action[0][MaskingRuleActionProfileId] = v.MaskingAction.MaskingProfileId
		rules[i][MaskingRuleAction] = action

		criteria := make([]map[string]interface{}, 1)
		criteria[0] = make(map[string]interface{})
		criteria[0][CriteriaCondition] = v.DataFilterCriteria.Condition

		identity := make([]map[string]interface{}, 1)
		identity[0] = *dataAccessIdentityToResource(&v.DataFilterCriteria.Identity)
		criteria[0][CriteriaIdentity] = identity

		rules[i][RuleCriteria] = criteria

		if v.ConditionalMasking.WhereCondition != "" {
			conditionalMasking := make([]map[string]interface{}, 1)
			conditionalMasking[0] = make(map[string]interface{})
			conditionalMasking[0][MaskingRuleWhereCondition] = v.ConditionalMasking.WhereCondition
			rules[i][MaskingRuleConditionMasking] = conditionalMasking
		}

	}
	out[0][MaskingRule] = rules
	return out
}

// //////////////////////////////////
// APIs
// //////////////////////////////////
func resourceSecurityPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	securityPolicyOutput, err, statusCode := c.GetSecurityPolicy(d.Id())
	if statusCode == 404 {
		d.SetId("")
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(SecurityPolicyName, securityPolicyOutput.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(SecurityPolicyProfile, securityProfilesToResource(securityPolicyOutput.SecurityProfiles, d)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceSecurityPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	input, err := resourceToSecurityPolicy(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = c.UpdateSecurityPolicy(d.Id(), input)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceSecurityPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	if err := c.DeleteSecurityPolicy(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
