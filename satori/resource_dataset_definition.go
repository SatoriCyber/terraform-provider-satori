package satori

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

func resourceDataSetDefinition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataSetDefinitionCreate,
		ReadContext:   resourceDataSetDefinitionRead,
		UpdateContext: resourceDataSetDefinitionUpdate,
		DeleteContext: resourceDataSetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Dataset definition configuration only.",
		Schema: map[string]*schema.Schema{
			"data_policy_id": getDatasetDataPolicyIdSchema(),
			"definition":     getDatasetDefinitionSchema(),
		},
	}
}

func resourceDataSetDefinitionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	if _, err := createDataSet(d, c); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDataSetDefinitionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	if _, err := getDataSet(c, d); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDataSetDefinitionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	if _, err := updateDataSet(d, c); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
