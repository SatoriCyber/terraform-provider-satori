package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
)

var (
	Name                        = "name"
	Hostname                    = "hostname"
	Id                          = "id"
	DataAccessControllerId      = "dataaccess_controller_id"
	CustomIngressPort           = "custom_ingress_port"
	OriginPort                  = "origin_port"
	ProjectIds                  = "project_ids"
	BaselineSecurityPolicy      = "baseline_security_policy"
	Type                        = "type"
	UnassociatedQueriesCategory = "unassociated_queries_category"
	UnsupportedQueriesCategory  = "unsupported_queries_category"
	Pattern                     = "pattern"
	ExcludedIdentities          = "excluded_identities"
	Exclusions                  = "exclusions"
	QueryAction                 = "query_action"
	ExcludedQueryPatterns       = "excluded_query_patterns"
	Identity                    = "identity"
	IdentityType                = "identity_type"
	NetworkPolicy               = "network_policy"
	AllowedRules                = "allowed_rules"
	BlockedRules                = "blocked_rules"
	Note                        = "note"
	IpRanges                    = "ip_ranges"
	IpRange                     = "ip_range"
)
var TreatAsMap = map[string]bool{
	Exclusions:                  true,
	UnsupportedQueriesCategory:  true,
	UnassociatedQueriesCategory: true,
	BaselineSecurityPolicy:      true,
	NetworkPolicy:               true,
}

func getDataStoreDefinitionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		Id: &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "DataStore resource id.",
		},
		Name: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "DataStore name.",
		},
		Hostname: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Data provider's FQDN hostname.", // example: snowflakecomputing.com, xyz.redshift.amazonaws.com:5439/dev
		},
		OriginPort: &schema.Schema{
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Port number description.",
		},
		DataAccessControllerId: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Host FQDN name.",
		},
		CustomIngressPort: &schema.Schema{
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     nil,
			Description: "Custom ingress port number description.",
		},
		ProjectIds: &schema.Schema{
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "ProjectIds list of project IDs",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		Type: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "IDs of Satori users that will be set as DataStore owners.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		BaselineSecurityPolicy: GetBaseLinePolicyDefinition(),
		NetworkPolicy:          GetNetworkPolicyDefinition(),
	}
}
func createDataStore(d *schema.ResourceData, c *api.Client) (*api.DataStoreOutput, error) {
	input, err := resourceToDataStore(d)

	if err != nil {
		return nil, err
	}

	result, err := c.CreateDataStore(input)
	if err != nil {
		return nil, err
	}
	d.SetId(result.Id)

	if err := d.Set("id", result.Id); err != nil {
		return nil, err
	}
	return result, err
}

// convert terraform resource defs into datastore type //
func resourceToDataStore(d *schema.ResourceData) (*api.DataStore, error) {
	out := api.DataStore{}

	baselineSecurityPolicyToResource, err := BaselineSecurityPolicyToResource(d.Get("baseline_security_policy").([]interface{}))
	if err != nil {
		return nil, err
	}

	networkPolicyToResource, err := NetworkPolicyToResource(d.Get(NetworkPolicy).([]interface{}))
	if err != nil {
		return nil, err
	}

	out.Name = d.Get("name").(string)
	out.Hostname = d.Get("hostname").(string)
	out.OriginPort = d.Get(OriginPort).(int)
	out.CustomIngressPort = d.Get("custom_ingress_port").(int)
	out.DataAccessControllerId = d.Get("dataaccess_controller_id").(string)
	out.ProjectIds = convertStringSet(d.Get("project_ids").(*schema.Set))
	out.Type = d.Get("type").(string)
	out.BaselineSecurityPolicy = baselineSecurityPolicyToResource
	out.NetworkPolicy = networkPolicyToResource
	return &out, nil
}

// update datastoreoutput
func getDataStore(c *api.Client, d *schema.ResourceData) (*api.DataStoreOutput, error) {
	result, err, statusCode := c.GetDataStore(d.Id())
	if statusCode == 404 {
		d.SetId("")
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	d.Set(Id, result.Id)
	d.Set(Name, result.Name)
	d.Set(Hostname, result.Hostname)
	d.Set(Type, result.Type)
	d.Set(OriginPort, result.OriginPort)
	d.Set(CustomIngressPort, result.CustomIngressPort)
	d.Set(DataAccessControllerId, result.DataAccessControllerId)
	d.Set(ProjectIds, result.ProjectIds)

	tfMap, err := GetBaseLinePolicyDatastoreOutput(result, err)
	if err != nil {
		return nil, err
	}
	d.Set(BaselineSecurityPolicy, []map[string]interface{}{tfMap})

	npMap, err := GetNetworkPolicyDatastoreOutput(result, err)
	if err != nil {
		return nil, err
	}
	d.Set(NetworkPolicy, []map[string]interface{}{npMap})

	return result, err
}

func extractMapFromInterface(in []interface{}) map[string]interface{} {
	if len(in) > 0 {
		if in[0] == nil {
			return nil
		}
		return in[0].(map[string]interface{})
	} else {
		return nil
	}
}

func updateDataStore(d *schema.ResourceData, c *api.Client) (*api.DataStoreOutput, error) {
	input, err := resourceToDataStore(d)
	if err != nil {
		return nil, err
	}
	result, err := c.UpdateDataStore(d.Id(), input)
	return result, err
}

func resourceDataStoreDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*api.Client)
	if err := c.DeleteDataStore(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return diags
}
