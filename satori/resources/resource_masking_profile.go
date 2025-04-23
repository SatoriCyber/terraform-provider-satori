package resources

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"log"
	"regexp"
)

func ResourceMaskingProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMaskingProfileCreate,
		ReadContext:   resourceMaskingProfileRead,
		UpdateContext: resourceMaskingProfileUpdate,
		DeleteContext: resourceMaskingProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Masking Profile.",
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Masking profile name.",
				ValidateFunc: func(v interface{}, key string) (warns []string, errs []error) {
					name := v.(string)
					var isValid = regexp.MustCompile(`^[a-zA-Z0-9 ]+$`).MatchString
					if !isValid(name) {
						errs = append(errs, fmt.Errorf("%q must contain only alphanumeric characters and spaces but got: %q", key, name))
					}
					return
				},
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Masking profile description.",
			},
			"condition": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "Masking profile condition.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type. Can be one of [TRUNCATE, TRUNCATE_END, REPLACE_CHAR, REPLACE_STRING, HASH, EMAIL_PREFIX, EMAIL_SUFFIX, EMAIL_FULL, EMAIL_HASH, CREDIT_CARD_PREFIX, CREDIT_CARD_FULL, CREDIT_CARD_HASH, IP_SUFFIX, IP_FULL, IP_HASH, DATE_YEAR_ONLY, DATE_1970_AGAIN, NO_ACTION, REDACT, NUMBER_ZERO, NUMBER_ROUND, ...]",
						},
						"replacement": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Replacement, relevant for: REPLACE_CHAR, REPLACE_STRING.",
						},
						"truncate": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Truncate, relevant for: TRUNCATE, TRUNCATE_END.",
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
func resourceToMaskingConditions(d *schema.ResourceData) *[]api.MaskingCondition {
	var maskConfigs []api.MaskingCondition
	if v, ok := d.GetOk("condition"); ok {
		maskConfigs = make([]api.MaskingCondition, len(v.([]interface{})))
		for i, raw := range v.([]interface{}) {
			inElement := raw.(map[string]interface{})
			maskConfigs[i] = api.MaskingCondition{}
			maskConfigs[i].Tag = inElement["tag"].(string)
			maskConfigs[i].Type = inElement["type"].(string)
			if inElement["truncate"].(int) > 0 {
				maskConfigs[i].Truncate = inElement["truncate"].(int)
			}
			replacement := inElement["replacement"].(string)
			maskConfigs[i].Replacement = &replacement
		}
	} else {
		maskConfigs = make([]api.MaskingCondition, 0)
	}
	return &maskConfigs
}

func resourceToMaskingProfile(d *schema.ResourceData) *api.MaskingProfile {
	maskingProfile := api.MaskingProfile{}
	maskingProfile.Name = d.Get("name").(string)
	description := d.Get("description").(string)
	maskingProfile.Description = &description
	maskingProfile.MaskConfigs = *resourceToMaskingConditions(d)

	return &maskingProfile
}

func resourceMaskingProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	input := resourceToMaskingProfile(d)

	result, err := c.CreateMaskingProfile(input)
	if err != nil {
		log.Printf("Recieved error in masking profile create: %s", err)
		return diag.FromErr(err)
	} else {
		d.SetId(result.Id)
	}

	return diags
}

// //////////////////////////////////
// Schema to resource mappers
// //////////////////////////////////
func maskingConditionToResource(conditions []api.MaskingCondition) []map[string]interface{} {
	out := make([]map[string]interface{}, len(conditions))
	for i, v := range conditions {
		out[i] = make(map[string]interface{})
		out[i]["type"] = v.Type
		out[i]["tag"] = v.Tag
		if v.Replacement != nil {
			out[i]["replacement"] = *v.Replacement
		}
		if v.Truncate > 0 {
			out[i]["truncate"] = v.Truncate
		}
	}
	return out
}

////////////////////////////////////
// APIs
////////////////////////////////////

func resourceMaskingProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	result, err, statusCode := c.GetMaskingProfile(d.Id())
	if statusCode == 404 {
		d.SetId("")
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", result.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := setNullableStringProp(result.Description, "description", d); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("condition", maskingConditionToResource(result.MaskConfigs)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceMaskingProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	input := resourceToMaskingProfile(d)
	_, err := c.UpdateMaskingProfile(d.Id(), input)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceMaskingProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	if err := c.DeleteMaskingProfile(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
