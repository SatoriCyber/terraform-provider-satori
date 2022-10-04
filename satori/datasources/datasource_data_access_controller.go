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
		Description: "Find DAC by type, region and cloud provider.",
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DAC's type. The available values are: PRIVATE, PRIVATE_MANAGED or PUBLIC.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DAC's region.",
			},
			"cloud_provider": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DAC's cloud provider.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DAC's ID.",
			},
			"ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "DAC's IPs list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	id := d.Get("id").(string)

	if id != "" {
		d.SetId(id)
		dac, err := c.QueryDataAccessControllerById(&id)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("region", dac.Region)
		if err != nil {
			return nil
		}
		err = d.Set("cloud_provider", dac.CloudProvider)
		if err != nil {
			return nil
		}
		if len(dac.Ips) > 0 {
			err := d.Set("ips", dac.Ips)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		return diags
	}

	if dacType == "PRIVATE" || dacType == "PRIVATE_MANAGED" {
		if id == "" {
			return diag.Errorf("DAC with type of '%s' must include its ID.", dacType)
		}
	}

	if region == "" || cloudProvider == "" {
		return diag.Errorf("Public DAC must include both region and cloud provider.")
	}

	dacs, err := c.QueryDataAccessControllers(&dacType, &region, &cloudProvider)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(dacs) == 0 {
		return diag.Errorf("No DAC found with type '%s', region '%s', and cloud provider '%s'!", dacType, region, cloudProvider)
	}

	if len(dacs) > 1 {
		return diag.Errorf("Got more than one DAC with the values: type '%s', region '%s' and cloud provider '%s'.", dacType, region, cloudProvider)
	}

	dac := dacs[0]
	d.SetId(dac.Id)
	if len(dac.Ips) > 0 {
		err := d.Set("ips", dac.Ips)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
