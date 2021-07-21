package satori

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

func resourceDirectoryGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDirectoryGroupCreate,
		ReadContext:   resourceDirectoryGroupRead,
		UpdateContext: resourceDirectoryGroupUpdate,
		DeleteContext: resourceDirectoryGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Directory group configuration.",
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group name.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group description.",
			},
			"member": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "Group members.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Member type, valid types are: USERNAME, IDP_GROUP, DB_ROLE, DIRECTORY_GROUP.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Member name for types: USERNAME, IDP_GROUP and DB_ROLE.",
						},
						"group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Directory group ID for member of type DIRECTORY_GROUP.",
						},
						"identity_provider": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Identity provider type for member of type IDP_GROUP, valid identity providers are: OKTA, AZURE, ONELOGIN",
						},
						"data_store_type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data store type for member of type DB_ROLE, valid types are: SNOWFLAKE, REDSHIFT, BIGQUERY, POSTGRESQL, ATHENA, MSSQL, SYNAPSE",
						},
					},
				},
			},
		},
	}
}

func resourceDirectoryGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	input := resourceToDirectoryGroup(d)

	result, err := c.CreateDirectoryGroup(input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*result.Id)
	return diags
}

func resourceToDirectoryGroup(d *schema.ResourceData) *api.DirectoryGroup {
	out := api.DirectoryGroup{}
	out.Name = d.Get("name").(string)
	if v, ok := d.GetOk("description"); ok {
		description := v.(string)
		out.Description = &description
	}
	if raw, ok := d.GetOk("member"); ok {
		in := raw.([]interface{})
		members := make([]api.DirectoryGroupMember, len(in))
		for i := 0; i < len(in); i++ {
			if v, ok := d.GetOk(fmt.Sprintf("member.%d.type", i)); ok {
				members[i].Type = v.(string)
			}
			if v, ok := d.GetOk(fmt.Sprintf("member.%d.name", i)); ok {
				name := v.(string)
				members[i].Name = &name
			}
			if v, ok := d.GetOk(fmt.Sprintf("member.%d.identity_provider", i)); ok {
				provider := v.(string)
				members[i].Provider = &provider
			}
			if v, ok := d.GetOk(fmt.Sprintf("member.%d.group_id", i)); ok {
				groupId := v.(string)
				members[i].GroupId = &groupId
			}
			if v, ok := d.GetOk(fmt.Sprintf("member.%d.data_store_type", i)); ok {
				dsType := v.(string)
				members[i].DsType = &dsType
			}
		}
		out.Members = members
	}
	return &out
}

func resourceDirectoryGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	result, err, statusCode := c.GetDirectoryGroup(d.Id())
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
	if result.Description != nil {
		if err := d.Set("description", result.Description); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if v, ok := d.GetOk("description"); ok && len(v.(string)) > 0 {
			if err := d.Set("description", nil); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if err := d.Set("member", directoryGroupMembersToResource(result.Members)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func directoryGroupMembersToResource(members []api.DirectoryGroupMember) []map[string]interface{} {
	out := make([]map[string]interface{}, len(members))
	for i, v := range members {
		out[i] = make(map[string]interface{})
		out[i]["type"] = v.Type
		if v.Name != nil && v.Type != "DIRECTORY_GROUP" {
			out[i]["name"] = *v.Name
		}
		if v.Provider != nil {
			out[i]["identity_provider"] = *v.Provider
		}
		if v.GroupId != nil {
			out[i]["group_id"] = *v.GroupId
		}
		if v.DsType != nil {
			out[i]["data_store_type"] = *v.DsType
		}
	}
	return out
}

func resourceDirectoryGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	input := resourceToDirectoryGroup(d)
	if _, err := c.UpdateDirectoryGroup(d.Id(), input); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDirectoryGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	err := c.DeleteDirectoryGroup(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
