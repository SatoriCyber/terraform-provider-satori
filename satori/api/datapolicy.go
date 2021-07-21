package api

const CustomPolicyDefaultPriority int = 100
const DataPolicyApiPrefix = "/api/data-policy"
const DataPolicyRulesPrefix = DataPolicyApiPrefix + "/rules"
const DataPolicySecurityPoliciesPrefix = DataPolicyApiPrefix + "/security-policies"
const DataPolicyPermissionsPrefix = DataPolicyApiPrefix + "/permissions"

type AccessControl struct {
	AccessControlEnabled bool `json:"permissionsEnabled"`
	UserRequestsEnabled  bool `json:"accessRequestsEnabled"`
	SelfServiceEnabled   bool `json:"selfServiceAccessEnabled"`
}

type CustomPolicy struct {
	Priority  int    `json:"priority"`
	RulesYaml string `json:"rulesYaml"`
	TagsYaml  string `json:"tagsYaml"`
}

type SecurityPolicies struct {
	Ids []string `json:"defaultSecurityPolicyIds"`
}

func (c *Client) UpdateCustomPolicy(id string, input *CustomPolicy) (*CustomPolicy, error) {
	output := CustomPolicy{}
	return &output, c.putJson(DataPolicyRulesPrefix, id, input, &output)
}

func (c *Client) GetCustomPolicy(id string) (*CustomPolicy, error) {
	output := CustomPolicy{}
	err, _ := c.getJsonById(DataPolicyRulesPrefix, id, &output)
	return &output, err
}

func (c *Client) UpdateSecurityPolicies(id string, input *SecurityPolicies) (*SecurityPolicies, error) {
	output := SecurityPolicies{}
	return &output, c.putJson(DataPolicySecurityPoliciesPrefix, id, input, &output)
}

func (c *Client) GetSecurityPolicies(id string) (*SecurityPolicies, error) {
	output := SecurityPolicies{}
	err, _ := c.getJsonById(DataPolicySecurityPoliciesPrefix, id, &output)
	return &output, err
}

func (c *Client) UpdateAccessControl(id string, input *AccessControl) (*AccessControl, error) {
	output := AccessControl{}
	return &output, c.putJson(DataPolicyPermissionsPrefix, id, input, &output)
}

func (c *Client) GetAccessControl(id string) (*AccessControl, error) {
	output := AccessControl{}
	err, _ := c.getJsonById(DataPolicyPermissionsPrefix, id, &output)
	return &output, err
}
