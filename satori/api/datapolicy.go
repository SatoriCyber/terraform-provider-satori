package api

const CustomPolicyDefaultPriority int = 100

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
	return &output, c.putJson("/api/data-policy/rules", id, input, &output)
}

func (c *Client) GetCustomPolicy(id string) (*CustomPolicy, error) {
	output := CustomPolicy{}
	return &output, c.getJson("/api/data-policy/rules", id, &output)
}

func (c *Client) UpdateSecurityPolicies(id string, input *SecurityPolicies) (*SecurityPolicies, error) {
	output := SecurityPolicies{}
	return &output, c.putJson("/api/data-policy/security-policies", id, input, &output)
}

func (c *Client) GetSecurityPolicies(id string) (*SecurityPolicies, error) {
	output := SecurityPolicies{}
	return &output, c.getJson("/api/data-policy/security-policies", id, &output)
}

func (c *Client) UpdateAccessControl(id string, input *AccessControl) (*AccessControl, error) {
	output := AccessControl{}
	return &output, c.putJson("/api/data-policy/permissions", id, input, &output)
}

func (c *Client) GetAccessControl(id string) (*AccessControl, error) {
	output := AccessControl{}
	return &output, c.getJson("/api/data-policy/permissions", id, &output)
}
