package satori

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"github.com/satoricyber/terraform-provider-satori/satori/datasources"
	"github.com/satoricyber/terraform-provider-satori/satori/resources"
	"strings"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		if s.Deprecated != "" {
			desc += " " + s.Deprecated
		}
		return strings.TrimSpace(desc)
	}
}

func NewProvider(version string) *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"satori_account": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Your Satori account ID.",
			},
			"service_account": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SATORI_SA", nil),
				Description: "Service account ID with administrative privileges." +
					" Can be specified with the `SATORI_SA` environment variable.",
			},
			"service_account_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SATORI_SA_KEY", nil),
				Description: "Service account key." +
					" Can be specified with the `SATORI_SA_KEY` environment variable.",
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
			"satori_dataset":                  resources.ResourceDataSet(),
			"satori_dataset_definition":       resources.ResourceDataSetDefinition(),
			"satori_directory_group":          resources.ResourceDirectoryGroup(),
			"satori_access_rule":              resources.ResourceDataAccessPermission(),
			"satori_self_service_access_rule": resources.ResourceDataAccessSelfServiceRule(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"satori_user": datasources.DatasourceUser(),
		},
	}

	p.ConfigureContextFunc = providerConfigure(version, p)

	return p
}

func providerConfigure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		username := d.Get("service_account").(string)
		password := d.Get("service_account_key").(string)
		verifyTls := d.Get("verify_tls").(bool)
		url := d.Get("url").(string)
		accountId := d.Get("satori_account").(string)

		userAgent := p.UserAgent("terraform-provider-satori", version)

		var diags diag.Diagnostics

		c, err := api.NewClient(&url, &userAgent, &accountId, &username, &password, verifyTls)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create Satori client",
				Detail:   "Unable to authenticate user",
			})
			return nil, diag.FromErr(err)
		}

		return c, diags
	}
}
