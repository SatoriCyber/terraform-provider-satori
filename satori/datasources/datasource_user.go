package datasources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"strings"
)

func DatasourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceUserRead,
		Description: "Find user ID by email",
		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "User's email address.",
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User's ID.",
			},
		},
	}
}

func datasourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	email := d.Get("email").(string)

	users, err := c.QueryUsers(&email)
	if err != nil {
		return diag.FromErr(err)
	}

	found := false
	if len(users) > 0 {
		for _, user := range users {
			if strings.EqualFold(user.Email, email) {
				d.SetId(user.Id)
				found = true
				break
			}
		}
	}

	if !found {
		return diag.Errorf("No user with email '%s' found!", email)
	}

	return diags
}
