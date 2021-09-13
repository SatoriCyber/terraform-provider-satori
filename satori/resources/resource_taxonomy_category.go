package resources

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"strings"
)

func ResourceTaxonomyCategory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCategoryCreate,
		ReadContext:   resourceCategoryRead,
		UpdateContext: resourceCategoryUpdate,
		DeleteContext: resourceCategoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Custom taxonomy category.",
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Category name.",
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
				Description: "Category tag.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Category description.",
			},
			"parent_category_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Parent category ID.",
			},
			"color": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Category color.",
			},
		},
	}
}

func resourceToCategory(d *schema.ResourceData) *api.TaxonomyCategory {
	out := api.TaxonomyCategory{}
	out.Name = d.Get("name").(string)
	out.Color = d.Get("color").(string)
	if v, ok := d.GetOk("description"); ok {
		description := v.(string)
		out.Description = &description
	}
	if v, ok := d.GetOk("parent_category_id"); ok {
		parentNode := v.(string)
		out.ParentNode = &parentNode
	}
	return &out
}

func resourceCategoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	input := resourceToCategory(d)

	result, err := c.CreateTaxonomyCategory(input)
	if err != nil {
		diag.FromErr(err)
	}

	d.SetId(result.Id)
	if err := d.Set("tag", result.Tag); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCategoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	result, err, statusCode := c.GetTaxonomyCategory(d.Id())
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
	if err := SetNullableStringProp(result.Description, "description", d); err != nil {
		return diag.FromErr(err)
	}
	if err := SetNullableStringProp(result.ParentNode, "parent_category_id", d); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("color", result.Color); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tag", result.Tag); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCategoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	input := resourceToCategory(d)
	result, err := c.UpdateTaxonomyCategory(d.Id(), input)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tag", result.Tag); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCategoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	if err := c.DeleteTaxonomyNode(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
