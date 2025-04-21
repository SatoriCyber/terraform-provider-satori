package resources

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"strings"
)

func ResourceTaxonomyClassifier() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClassifierCreate,
		ReadContext:   resourceClassifierRead,
		UpdateContext: resourceClassifierUpdate,
		DeleteContext: resourceClassifierDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Custom taxonomy classifier.",
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Classifier name.",
				ValidateFunc: func(v interface{}, key string) (warns []string, errs []error) {
					name := v.(string)
					if strings.Contains(name, ".") || strings.Contains(name, ":") {
						errs = append(errs, fmt.Errorf("%q must not include '.' or ':' but got: %q", key, name))
					}
					return
				},
			},
			"tag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Classifier tag.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Classifier description.",
			},
			"parent_category": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Parent category ID.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Classifier type, valid types are: NON_AUTOMATIC, CUSTOM, SATORI_BASED.",
				ValidateFunc: func(v interface{}, key string) (warns []string, errs []error) {
					value := v.(string)
					if value != "NON_AUTOMATIC" && value != "CUSTOM" && value != "SATORI_BASED" {
						errs = append(errs, fmt.Errorf("%q must be one of 'NON_AUTOMATIC, CUSTOM or SATORI_BASED' but got: %q", key, value))
					}
					return
				},
			},
			"custom_config": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Configuration for CUSTOM classifier type.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_name_pattern": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Field name pattern.",
						},
						"field_type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field type, valid types are: ANY, TEXT, NUMERIC, DATE.",
						},
						"values": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of possible field values.",
							ConflictsWith: []string{
								"custom_config.0.value_pattern",
							},
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"value_pattern": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Value pattern.",
							ConflictsWith: []string{
								"custom_config.0.values",
							},
						},
						"value_case_sensitive": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Should value match be case sensitive.",
							Default:     true,
						},
					},
				},
			},
			"satori_based_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Configuration for SATORI_BASED classifier type.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"satori_base_classifier": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Base Satori classifier ID.\nSee https://satoricyber.com/docs/taxonomy/standard-classifiers for a list of possible values.",
						},
					},
				},
			},
			"additional_satori_categories": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of additional Satori taxonomy category IDs.\nSee https://satoricyber.com/docs/taxonomy/standard-categories for a list of possible values.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"scope": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Scope of relevant locations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datasets": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "IDs of datasets to include in the scope.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"include_location": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Location to include in the scope.",
							Elem:        getDatasetLocationResource(),
						},
					},
				},
			},
		},
	}
}

func resourceToClassifier(d *schema.ResourceData) (*api.TaxonomyClassifier, error) {
	out := api.TaxonomyClassifier{}
	out.Name = d.Get("name").(string)
	out.ParentNode = d.Get("parent_category").(string)
	if v, ok := d.GetOk("description"); ok {
		description := v.(string)
		out.Description = &description
	}

	out.Config.Type = d.Get("type").(string)

	if out.Config.Type == "SATORI_BASED" {
		if v, ok := d.GetOk("satori_based_config.0.satori_base_classifier"); ok {
			value := v.(string)
			out.Config.SatoriBaseClassifierId = &value
		}
	} else if out.Config.Type == "CUSTOM" {
		if v, ok := d.GetOk("custom_config.0.field_type"); ok {
			value := v.(string)
			out.Config.FieldType = &value
		}
		if v, ok := d.GetOk("custom_config.0.field_name_pattern"); ok {
			value := v.(string)
			out.Config.FieldNamePattern = &value
		}
		caseInsensitive := false
		if v, ok := d.GetOk("custom_config.0.value_case_sensitive"); ok {
			caseInsensitive = !v.(bool)
		}
		if _, ok := d.GetOk("custom_config.0.values"); ok {
			var values api.TaxonomyClassifierValues
			valuesConfig, _ := getStringListProp("custom_config.0.values", d)
			values.Values = valuesConfig
			values.CaseInsensitive = caseInsensitive
			out.Config.Values = &values
		} else if v, ok := d.GetOk("custom_config.0.value_pattern"); ok {
			value := []string{v.(string)}
			var values api.TaxonomyClassifierValues
			values.Values = &value
			values.Regex = true
			values.CaseInsensitive = caseInsensitive
			out.Config.Values = &values
		}
	}

	additionalSatoiCategories, _ := getStringListProp("additional_satori_categories", d)
	out.Config.AdditionalCategories = *additionalSatoiCategories

	scopeDataset, _ := getStringListProp("scope.0.datasets", d)
	out.Scope.DatasetIds = *scopeDataset
	locationOutput, err := resourceToLocations(d, "scope.0.include_location", false)
	if err != nil {
		return nil, err
	}
	out.Scope.IncludeLocations = *locationOutput

	return &out, nil
}

func resourceClassifierCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	input, err := resourceToClassifier(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := c.CreateTaxonomyClassifier(input)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.Id)
	if err := d.Set("tag", result.Tag); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceClassifierRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	result, err, statusCode := c.GetTaxonomyClassifier(d.Id())
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
	if err := d.Set("parent_category", result.ParentNode); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tag", result.Tag); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", result.Config.Type); err != nil {
		return diag.FromErr(err)
	}

	satoriConfig := make(map[string]interface{})
	if result.Config.SatoriBaseClassifierId != nil {
		satoriConfig["satori_base_classifier"] = *result.Config.SatoriBaseClassifierId
	}
	if err := setMapProp(&satoriConfig, "satori_based_config", d); err != nil {
		return diag.FromErr(err)
	}

	customConfig := make(map[string]interface{})
	if result.Config.FieldNamePattern != nil {
		customConfig["field_name_pattern"] = *result.Config.FieldNamePattern
	}
	if result.Config.FieldType != nil {
		customConfig["field_type"] = *result.Config.FieldType
	}
	if result.Config.Values != nil && result.Config.Values.Values != nil {
		if result.Config.Values.Regex {
			customConfig["value_pattern"] = (*result.Config.Values.Values)[0]
		} else {
			customConfig["values"] = *result.Config.Values.Values
		}
		customConfig["value_case_sensitive"] = !result.Config.Values.CaseInsensitive
	}
	if err := setMapProp(&customConfig, "custom_config", d); err != nil {
		return diag.FromErr(err)
	}

	if err := setStringListProp(&result.Config.AdditionalCategories, "additional_satori_categories", d); err != nil {
		return diag.FromErr(err)
	}

	scope := make(map[string]interface{})
	if len(result.Scope.DatasetIds) > 0 {
		scope["datasets"] = result.Scope.DatasetIds
	}
	if len(result.Scope.IncludeLocations) > 0 {
		scope["include_location"] = locationsToResource(&result.Scope.IncludeLocations, d, "scope.0.include_location", Location)
	}
	if err := setMapProp(&scope, "scope", d); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceClassifierUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	input, err := resourceToClassifier(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := c.UpdateTaxonomyClassifier(d.Id(), input)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tag", result.Tag); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceClassifierDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	if err := c.DeleteTaxonomyNode(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
