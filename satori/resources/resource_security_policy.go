package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"log"
)

var (
	MaskingActive                      = "active"
	SecurityPolicyName                 = "name"
	SecurityPolicyProfile              = "profile"
	MaskingProfile                     = "masking"
	MaskingRule                        = "rule"
	MaskingRuleId                      = "id"
	MaskingRuleDescription             = "description"
	MaskingRuleActive                  = "active"
	MaskingRuleAction                  = "action"
	MaskingRuleActionType              = "type"
	MaskingRuleActionProfileId         = "masking_profile_id"
	MaskingRuleCriteria                = "criteria"
	MaskingRuleCriteriaCondition       = "condition"
	MaskingRuleCriteriaIdentity        = "identity"
	MaskingRuleActionDefaultActionType = "APPLY_MASKING_PROFILE"
	RowLevelSecurity                   = "row_level_security"
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
			SecurityPolicyName: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Security policy name.",
			},
			SecurityPolicyProfile: &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Security policy profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						MaskingProfile: resourceMaskingProfile(),
						//RowLevelSecurity: resourceRowLevelSecurity(),
					},
				},
			},
		},
	}
}

func resourceRowLevelSecurity() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Row level security, TBD",
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
				MaskingActive: &schema.Schema{
					Type:        schema.TypeBool,
					Required:    true,
					Description: "Is active.",
				},
				MaskingRule: &schema.Schema{
					Type:        schema.TypeList,
					Required:    true,
					Description: "Masking Rule.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							MaskingRuleId: &schema.Schema{
								Type:        schema.TypeString,
								Required:    true,
								Description: "Rule id, has to be unique.",
							},
							MaskingRuleDescription: &schema.Schema{
								Type:        schema.TypeString,
								Required:    true,
								Description: "Rule description.",
							},
							MaskingRuleActive: &schema.Schema{
								Type:        schema.TypeBool,
								Required:    true,
								Description: "Is active rule.",
							},
							MaskingRuleAction: &schema.Schema{
								Type:        schema.TypeList,
								Required:    true,
								MaxItems:    1,
								Description: "Rule action.",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										MaskingRuleActionType: &schema.Schema{
											Type:        schema.TypeString,
											Optional:    true,
											Default:     MaskingRuleActionDefaultActionType,
											Description: "Rule type.",
										},
										MaskingRuleActionProfileId: &schema.Schema{
											Type:        schema.TypeString,
											Required:    true,
											Description: "The reference id to be applied as masking profile.",
										},
									},
								},
							},
							MaskingRuleCriteria: &schema.Schema{
								Type:        schema.TypeList,
								Required:    true,
								MinItems:    1,
								MaxItems:    1,
								Description: "Masking criteria.",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										MaskingRuleCriteriaCondition: &schema.Schema{
											Type:        schema.TypeString,
											Required:    true,
											Description: "Identity condition, for example IS_NOT, IS, etc.",
										},
										MaskingRuleCriteriaIdentity: resourceDataAccessIdentity(),
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

////////////////////////////////////
// Resource to schema mappers
////////////////////////////////////
func resourceToSecurityProfiles(d *schema.ResourceData) *api.SecurityProfiles {
	out := api.SecurityProfiles{}
	if _, ok := d.GetOk(SecurityPolicyProfile); ok {
		out.Masking = api.MaskingSecurityProfile{}
		if m, ok := d.GetOk("profile.0.masking.0"); ok {
			resourceToMasking(d, m, &out)
		}
	}
	//if _, ok := d.GetOk("profile.0.row_level_security.0"); ok {
	//	out.RowLevelSecurity = *(new(api.RowLevelSecurityProfile))
	//}

	return &out
}

func resourceToMasking(d *schema.ResourceData, m interface{}, out *api.SecurityProfiles) {
	masking := m.(map[string]interface{})

	isActive := masking[MaskingActive].(bool)

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
	outElement.Id = inElement[MaskingRuleId].(string)
	outElement.Description = inElement[MaskingRuleDescription].(string)
	outElement.Active = inElement[MaskingActive].(bool)

	actionList := inElement[MaskingRuleAction].([]interface{})
	action := actionList[0].(map[string]interface{})
	log.Printf("Action: %s", action)
	maskingProfileId := action[MaskingRuleActionProfileId].(string)

	// Masking action
	outElement.MaskingAction.MaskingProfileId = maskingProfileId
	outElement.MaskingAction.Type = MaskingRuleActionDefaultActionType

	// Masking criteria
	criteriaList := inElement[MaskingRuleCriteria].([]interface{})
	criteria := criteriaList[0].(map[string]interface{})
	outElement.DataFilterCriteria.Condition = criteria[MaskingRuleCriteriaCondition].(string)

	identityList := criteria[MaskingRuleCriteriaIdentity].([]interface{})

	identity := identityList[0].(map[string]interface{})
	outElement.DataFilterCriteria.Identity = *resourceToIdentity(identity)
	(*rules)[i] = outElement
}

func resourceToSecurityPolicy(d *schema.ResourceData) *api.SecurityPolicy {
	out := api.SecurityPolicy{}
	out.Name = d.Get(SecurityPolicyName).(string)

	out.SecurityProfiles = *resourceToSecurityProfiles(d)
	return &out
}

func resourceSecurityPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	input := resourceToSecurityPolicy(d)

	result, err := c.CreateSecurityPolicy(input)
	if err != nil {
		log.Printf("Recieved error in masking profile create: %s", err)
		diag.FromErr(err)
	} else {
		d.SetId(result.Id)
	}

	return diags
}

////////////////////////////////////
// Schema to resource mappers
////////////////////////////////////
func securityProfilesToResource(profiles api.SecurityProfiles) interface{} {

	out := make([]map[string]interface{}, 1)

	out[0] = make(map[string]interface{})
	out[0][MaskingProfile] = maskingToResource(profiles.Masking)
	//out[0][RowLevelSecurity] = rowLevelSecurityToResource(profiles.RowLevelSecurity)

	return out
}

func rowLevelSecurityToResource(security api.RowLevelSecurityProfile) interface{} {
	//out := make([]map[string]interface{}, 0)

	return ""
}

func maskingToResource(masking api.MaskingSecurityProfile) interface{} {
	out := make([]map[string]interface{}, 1)
	out[0] = make(map[string]interface{})

	out[0] = make(map[string]interface{})
	out[0][MaskingActive] = masking.Active
	rules := make([]map[string]interface{}, len(masking.Rules))

	for i, v := range masking.Rules {
		rules[i] = make(map[string]interface{})
		rules[i][MaskingRuleActive] = v.Active
		rules[i][MaskingRuleDescription] = v.Description
		rules[i][MaskingRuleId] = v.Id

		action := make([]map[string]interface{}, 1)
		action[0] = make(map[string]interface{})
		action[0][MaskingRuleActionType] = v.MaskingAction.Type
		action[0][MaskingRuleActionProfileId] = v.MaskingAction.MaskingProfileId
		rules[i][MaskingRuleAction] = action

		criteria := make([]map[string]interface{}, 1)
		criteria[0] = make(map[string]interface{})
		criteria[0][MaskingRuleCriteriaCondition] = v.DataFilterCriteria.Condition

		identity := make([]map[string]interface{}, 1)
		identity[0] = *dataAccessIdentityToResource(&v.DataFilterCriteria.Identity)
		criteria[0][MaskingRuleCriteriaIdentity] = identity
		rules[i][MaskingRuleCriteria] = criteria
	}
	out[0][MaskingRule] = rules
	return out
}

////////////////////////////////////
// APIs
////////////////////////////////////
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

	if err := d.Set(SecurityPolicyProfile, securityProfilesToResource(securityPolicyOutput.SecurityProfiles)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceSecurityPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	input := resourceToSecurityPolicy(d)
	_, err := c.UpdateSecurityPolicy(d.Id(), input)
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
