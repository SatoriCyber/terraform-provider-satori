package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

func ResourceDataStoreDefinition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataStoreDefinitionCreate,
		ReadContext:   resourceDataStoreDefinitionRead,
		UpdateContext: resourceDataStoreDefinitionUpdate,
		DeleteContext: resourceDataStoreDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "DataStore definition configuration only.",
		Schema: map[string]*schema.Schema{
			"datastore_id": getDataStoreDataPolicyIdSchema(),
			"definition":   getDataStoreDefinitionSchema(),
			"hostname":     getDataStoreDefinitionSchema(),
			"port":         getDataStoreDefinitionSchema(),
		},
	}
}

func resourceDataStoreDefinitionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	if _, err := createDataStore(d, c); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDataStoreDefinitionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	if _, err := getDataStore(c, d); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDataStoreDefinitionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	if _, err := updateDataStore(d, c); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
