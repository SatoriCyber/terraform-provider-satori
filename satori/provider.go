package satori

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"service_account_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SATORI_SA_ID", nil),
			},
			"service_account_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SATORI_SA_KEY", nil),
			},
			"verify_tls": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.HostURL,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"satori_dataset": resourceDataSet(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("service_account_id").(string)
	password := d.Get("service_account_key").(string)
	verifyTls := d.Get("verify_tls").(bool)
	url := d.Get("url").(string)
	accountId := d.Get("account_id").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c, err := api.NewClient(&url, &accountId, &username, &password, verifyTls)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Satori client",
			Detail:   "Unable to auth user for authenticated client",
		})
		return nil, diag.FromErr(err)
	}

	return c, diags
}
