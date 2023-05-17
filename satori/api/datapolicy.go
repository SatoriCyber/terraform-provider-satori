package api

const CustomPolicyDefaultPriority int = 100

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
