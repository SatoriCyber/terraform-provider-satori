package datasources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

func DatasourceDataAccessController() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceDataAccessControllerRead,
		Description: "Find DAC by type, region, cloud provider and unique name.",
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DAC's type. The available values are: PRIVATE, PRIVATE_MANAGED or PUBLIC.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DAC's region.",
			},
			"cloud_provider": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DAC's cloud provider.",
			},
			"unique_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DAC's unique name. Provide when the type is PRIVATE or PRIVATE_MANAGED.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DAC's ID.",
			},
		},
	}
}

func datasourceDataAccessControllerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	dacType := d.Get("type").(string)
	region := d.Get("region").(string)
	cloudProvider := d.Get("cloud_provider").(string)
	uniqueName := d.Get("unique_name").(string)

	if (dacType == "PRIVATE" || dacType == "PRIVATE_MANAGED") && uniqueName == "" {
		return diag.Errorf("DAC with type of '%s' must include its unique name. region: '%s', cloud provider: '%s'", dacType, region, cloudProvider)
	}

	dacs, err := c.QueryDataAccessControllers(&dacType, &region, &cloudProvider, &uniqueName)
	if err != nil {
		return diag.FromErr(err)
	}

	found := false
	if len(dacs) > 0 {
		for _, dac := range dacs {
			if dac.Type == dacType && dac.Region == region && dac.CloudProvider == cloudProvider && (uniqueName == "" || dac.UniqueName == uniqueName) {
				d.SetId(dac.Id)
				found = true
				break
			}
		}
	}

	if !found {
		return diag.Errorf("No DAC found with type '%s', region '%s', cloud provider '%s' and unique name '%s' !", dacType, region, cloudProvider, uniqueName)
	}

	return diags
}
