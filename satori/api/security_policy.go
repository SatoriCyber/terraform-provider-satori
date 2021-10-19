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

type SecurityPolicyOutput struct {
	SecurityPolicy
	Id string `json:"id"`
}

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

type RowLevelSecurityProfile struct {
	Active bool `json:"active"`
	//Maps
}

type SecurityProfiles struct {
	Masking          MaskingSecurityProfile  `json:"masking"`
	RowLevelSecurity RowLevelSecurityProfile `json:"rowLevelSecurity"`
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
