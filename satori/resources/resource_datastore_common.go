package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/satoricyber/terraform-provider-satori/satori/api"
	"log"
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
	DeploymentType              = "deployment_type"
	AwsHostedZoneId             = "aws_hosted_zone_id"
	AwsServerRoleArn            = "aws_service_role_arn"
	UnassociatedQueriesCategory = "unassociated_queries_category"
	Credentials                 = "credentials"
	Enabled                     = "enabled"
	Username                    = "username"
	Password                    = "password"
	UnsupportedQueriesCategory  = "unsupported_queries_category"
	Pattern                     = "pattern"
	ExcludedIdentities          = "excluded_identities"
	Exclusions                  = "exclusions"
	QueryAction                 = "query_action"
	ExcludedQueryPatterns       = "excluded_query_patterns"
	Identity                    = "identity"
	IdentityType                = "identity_type"
	NetworkPolicy               = "network_policy"
	SatoriAuthSettings          = "satori_auth_settings"
	DataStoreSettings           = "datastore_settings"
	DatabricksSettings          = "databricks_settings"
	AllowedRules                = "allowed_rules"
	BlockedRules                = "blocked_rules"
	Note                        = "note"
	IpRanges                    = "ip_ranges"
	IpRange                     = "ip_range"
	SatoriHostname              = "satori_hostname"
	EnablePersonalAccessToken   = "enable_personal_access_token"
	DatabricksAccountId         = "account_id"
	DatabricksWarehouseId       = "warehouse_id"
	DatabricksClientId          = "client_id"
	DatabricksClientSecret      = "client_secret"
)
var TreatAsMap = map[string]bool{
	Exclusions:                  true,
	UnsupportedQueriesCategory:  true,
	UnassociatedQueriesCategory: true,
	BaselineSecurityPolicy:      true,
	NetworkPolicy:               true,
	Credentials:                 true,
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
		SatoriHostname: &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Satori Hostname.",
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
			Description: "The datastore type, for example: POSTGRESQL, SNOWFLAKE, etc. The full list is available at https://app.satoricyber.com/docs/api#post-/api/v1/datastore",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		SatoriAuthSettings:     GetSatoriAuthSettingsDefinitions(),
		BaselineSecurityPolicy: GetBaseLinePolicyDefinition(),
		NetworkPolicy:          GetNetworkPolicyDefinition(),
		DataStoreSettings:      GetDataStoreSettingsDefinition(),
		DatabricksSettings:     GetDatabricksSettingsDefinition(),
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
	dataStoreType := d.Get("type").(string)

	baselineSecurityPolicyToResource, err := BaselineSecurityPolicyToResource(d.Get("baseline_security_policy").([]interface{}))
	if err != nil {
		return nil, err
	}

	networkPolicyToResource, err := NetworkPolicyToResource(d.Get(NetworkPolicy).([]interface{}))
	if err != nil {
		return nil, err
	}

	satoriAuthSettingsToResource, err := SatoriAuthSettingsToResource(d, d.Get(SatoriAuthSettings).([]interface{}))
	if err != nil {
		return nil, err
	}

	var dataStoreSettingsToResource *api.DataStoreSettings

	if dataStoreType == "DATABRICKS" {
		dataStoreSettingsToResource, err = DataStoreSettingsToResource(d.Get(DatabricksSettings).([]interface{}))
		if err != nil {
			return nil, err
		}
	} else {
		dataStoreSettingsToResource, err = DataStoreSettingsToResource(d.Get(DataStoreSettings).([]interface{}))
		if err != nil {
			return nil, err
		}
	}

	out.Name = d.Get("name").(string)
	out.Hostname = d.Get("hostname").(string)
	out.OriginPort = d.Get(OriginPort).(int)
	out.CustomIngressPort = d.Get("custom_ingress_port").(int)
	out.DataAccessControllerId = d.Get("dataaccess_controller_id").(string)
	out.ProjectIds = convertStringSet(d.Get("project_ids").(*schema.Set))

	out.Type = dataStoreType
	out.BaselineSecurityPolicy = baselineSecurityPolicyToResource
	out.NetworkPolicy = networkPolicyToResource
	out.SatoriAuthSettings = satoriAuthSettingsToResource
	out.DataStoreSettings = dataStoreSettingsToResource
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
	d.Set(SatoriHostname, result.SatoriHostname)
	datastoreType := result.Type
	d.Set(Type, datastoreType)
	d.Set(OriginPort, result.OriginPort)
	d.Set(CustomIngressPort, result.CustomIngressPort)
	d.Set(DataAccessControllerId, result.DataAccessControllerId)
	d.Set(ProjectIds, result.ProjectIds)

	tfMap, err := GetBaseLinePolicyDatastoreOutput(result, err)
	if err != nil {
		return nil, err
	}
	if len(tfMap) > 0 { // empty result, skip it.
		d.Set(BaselineSecurityPolicy, []map[string]interface{}{tfMap})
	}

	npMap, err := GetNetworkPolicyDatastoreOutput(result, err)
	if err != nil {
		return nil, err
	}
	if len(npMap) > 0 { // empty result, skip it.
		d.Set(NetworkPolicy, []map[string]interface{}{npMap})
	}

	sasMap, err := GetSatoriAuthSettingsDatastoreOutput(d, result, err)
	if err != nil {
		return nil, err
	}
	if len(sasMap) > 0 { // empty result, skip it.
		d.Set(SatoriAuthSettings, []map[string]interface{}{sasMap})
	}

	dsSettingsMap, err := GetSatoriDatastoreSettingsDatastoreOutput(result, err)
	if err != nil {
		return nil, err
	}
	if len(dsSettingsMap) > 0 { // empty result, skip it.
		if datastoreType == "DATABRICKS" {
			d.Set(DatabricksSettings, []map[string]interface{}{dsSettingsMap})
		} else {
			d.Set(DataStoreSettings, []map[string]interface{}{dsSettingsMap})

		}
	}
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
	if err != nil {
		// Handle the error and restore the configuration
		restoreConfiguration(d)
	}

	return result, err
}

// This is a tricky part of the terraform state management:
// - Returning an error diagnostic does not stop the state from being updated.
// - see more info here: https://developer.hashicorp.com/terraform/plugin/framework/diagnostics#how-errors-affect-state
// Therefore, in error scenario (like timeout or server error, and case the password is changed, we need to restore the configuration,
// otherwise the password will be not be detected as `changed` next time the `terraform plan` is called.
//
// At this point we have to restore the password to the old value only while all other properties are updated from backend response.
func restoreConfiguration(d *schema.ResourceData) {
	log.Printf("Failed to update the data store resource, restoring configuration in the state...")
	passwordResourcePath := "satori_auth_settings.0.credentials.0.password"
	if d.HasChange(passwordResourcePath) {
		oldV, _ := d.GetChange(passwordResourcePath)
		log.Printf("The password has changed from state, overriding it with the old value")
		satoriAuthSettingsInterface := d.Get("satori_auth_settings").([]interface{})
		//log.Printf("Current state of satori_auth_settings is: %v", satoriAuthSettingsInterface)
		if len(satoriAuthSettingsInterface) != 0 {
			satori_auth_setting := satoriAuthSettingsInterface[0].(map[string]interface{})
			credentialsMap := satori_auth_setting[Credentials].([]interface{})
			//log.Printf("Current state of credentialsMap is: %v", credentialsMap)
			if len(credentialsMap) > 0 { // found credentials object
				credentials := credentialsMap[0].(map[string]interface{})
				credentials[Password] = oldV.(string)
			}
			//log.Printf("After the change, state of satori_auth_settings is: %v", satoriAuthSettingsInterface)
			// push config to state
			if err := d.Set("satori_auth_settings", satoriAuthSettingsInterface); err != nil {
				log.Printf("Failed to set satori_auth_settings: %v", err)
			}
		}
	} else {
		log.Printf("The password hasn't change from state, no need to restore it")
	}
}

func resourceDataStoreDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*api.Client)
	if err := c.DeleteDataStore(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return diags
}
