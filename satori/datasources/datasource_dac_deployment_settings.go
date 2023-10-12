package datasources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

func DatasourceDeploymentSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceDeploymentSettingsRead,
		Description: "Get deployment settings by DAC id",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DAC's id.",
			},
			"service_account": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The DAC's deployments service account",
			},
		},
	}
}

func datasourceDeploymentSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	dacId := d.Get("id").(string)

	if dacId != "" {
		d.SetId(dacId)
		deploymentSettings, err := c.QueryDeploymentSettings(&dacId)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("service_account", deploymentSettings.GSA)
		if err != nil {
			return nil
		}
		return diags
	}

	return diags
}
