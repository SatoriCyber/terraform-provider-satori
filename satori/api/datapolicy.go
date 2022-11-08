package api

const CustomPolicyDefaultPriority int = 100
const DataPolicyApiPrefix = "/api/v1/data-policy"
const DataPolicyRulesSuffix = "rules"
const DataPolicySecurityPoliciesSuffix = "security-policies"
const DataPolicyPermissionsSuffix = "permissions"

type AccessControl struct {
	AccessControlEnabled bool `json:"permissionsEnabled"`
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
	return &output, c.putJson(DataPolicyApiPrefix, DataPolicyRulesSuffix, id, input, &output)
}

func (c *Client) GetCustomPolicy(id string) (*CustomPolicy, error) {
	output := CustomPolicy{}
	err, _ := c.getJsonById(DataPolicyApiPrefix, DataPolicyRulesSuffix, id, &output)
	return &output, err
}

func (c *Client) UpdateSecurityPolicies(id string, input *SecurityPolicies) (*SecurityPolicies, error) {
	output := SecurityPolicies{}
	return &output, c.putJson(DataPolicyApiPrefix, DataPolicySecurityPoliciesSuffix, id, input, &output)
}

func (c *Client) GetSecurityPolicies(id string) (*SecurityPolicies, error) {
	output := SecurityPolicies{}
	err, _ := c.getJsonById(DataPolicyApiPrefix, DataPolicySecurityPoliciesSuffix, id, &output)
	return &output, err
}

func (c *Client) UpdateAccessControl(id string, input *AccessControl) (*AccessControl, error) {
	output := AccessControl{}
	return &output, c.putJson(DataPolicyApiPrefix, DataPolicyPermissionsSuffix, id, input, &output)
}

func (c *Client) GetAccessControl(id string) (*AccessControl, error) {
	output := AccessControl{}
	err, _ := c.getJsonById(DataPolicyApiPrefix, DataPolicyPermissionsSuffix, id, &output)
	return &output, err
}
