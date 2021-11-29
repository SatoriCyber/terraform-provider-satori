package api

import (
	"encoding/json"
	"log"
)

const SecurityPolicyApiPrefix = "/api/v1/security-policies"

type SecurityPolicy struct {
	Name             string           `json:"name"`
	SecurityProfiles SecurityProfiles `json:"profiles"`
}

type SecurityProfiles struct {
	Masking          *MaskingSecurityProfile  `json:"masking,omitempty"`
	RowLevelSecurity *RowLevelSecurityProfile `json:"rowLevelSecurity,omitempty"`
}

type SecurityPolicyOutput struct {
	SecurityPolicy
	Id string `json:"id"`
}

/////////////////////
// Masking
/////////////////////
type MaskingSecurityProfile struct {
	Active bool          `json:"active"`
	Rules  []MaskingRule `json:"rules"`
}

type DataFilterCriteria struct {
	Condition string             `json:"condition"`
	Identity  DataAccessIdentity `json:"identity"`
}

type MaskingAction struct {
	MaskingProfileId string `json:"maskingProfileId"`
	Type             string `json:"type"`
}

type MaskingRule struct {
	Id                 string             `json:"id"`
	Active             bool               `json:"active"`
	Description        string             `json:"description"`
	DataFilterCriteria DataFilterCriteria `json:"criteria"`
	MaskingAction      MaskingAction      `json:"maskingAction"`
}

/////////////////////
// Row Level Security
/////////////////////
type RowLevelSecurityRuleFilter struct {
	DataStoreId    string                  `json:"dataStoreId"`
	LocationPrefix *DataSetGenericLocation `json:"locationPrefix"`
	LogicYaml      string                  `json:"logicYaml"`
	Advanced       bool                    `json:"advanced"`
}

type RowLevelSecurityRule struct {
	Id          string                     `json:"id"`
	Active      bool                       `json:"active"`
	Description string                     `json:"description"`
	RuleFilter  RowLevelSecurityRuleFilter `json:"filter"`
}

type RowLevelSecurityMapDataFilter struct {
	Criteria DataFilterCriteria `json:"criteria"`
	Values   DataFilterValues   `json:"values"`
}

type DataFilterValues struct {
	Type   string    `json:"type"` //   STRING, NUMERIC, ANY_VALUE, ALL_OTHER_VALUES
	Values *[]string `json:"values"`
}

type DataFilterDefaultValues struct {
	Type   string    `json:"type"` //   STRING, NUMERIC, NO_VALUE, ALL_OTHER_VALUES
	Values *[]string `json:"values,omitempty"`
}

type RowLevelSecurityFilter struct {
	Name     string                          `json:"name"`
	Filters  []RowLevelSecurityMapDataFilter `json:"filters"`
	Defaults DataFilterDefaultValues         `json:"defaults"`
}

type RowLevelSecurityProfile struct {
	Active bool                     `json:"active"`
	Rules  []RowLevelSecurityRule   `json:"rules"`
	Maps   []RowLevelSecurityFilter `json:"maps"`
}

func (c *Client) CreateSecurityPolicy(input *SecurityPolicy) (*SecurityPolicyOutput, error) {
	output := SecurityPolicyOutput{}
	jsonInput, _ := json.Marshal(input)
	log.Printf("Going to create security policy %s", jsonInput)
	return &output, c.postJsonForAccount(SecurityPolicyApiPrefix, input, &output)
}

func (c *Client) UpdateSecurityPolicy(id string, input *SecurityPolicy) (*SecurityPolicyOutput, error) {
	output := SecurityPolicyOutput{}
	return &output, c.putJson(SecurityPolicyApiPrefix, "", id, input, &output)
}

func (c *Client) GetSecurityPolicy(id string) (*SecurityPolicyOutput, error, int) {
	output := SecurityPolicyOutput{}
	err, statusCode := c.getJsonById(SecurityPolicyApiPrefix, "", id, &output)
	return &output, err, statusCode
}

func (c *Client) DeleteSecurityPolicy(id string) error {
	return c.delete(SecurityPolicyApiPrefix, id)
}
