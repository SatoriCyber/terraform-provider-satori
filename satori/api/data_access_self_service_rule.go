package api

import "strconv"

const DataAccessSelfServiceApiPrefix = "/api/v1/data-access-rule/self-service"

type DataAccessSelfServiceRule struct {
	Id               *string                                  `json:"id,omitempty"`
	ParentId         *string                                  `json:"parentId,omitempty"`
	Suspended        bool                                     `json:"suspended,omitempty"`
	Identity         *DataAccessIdentity                      `json:"identity,omitempty"`
	AccessLevel      string                                   `json:"accessLevel"`
	TimeLimit        DataAccessSelfServiceAndRequestTimeLimit `json:"timeLimit"`
	UnusedTimeLimit  DataAccessUnusedTimeLimit                `json:"unusedTimeLimit"`
	SecurityPolicies *[]string                                `json:"securityPolicyIds,omitempty"`
}

func (c *Client) CreateDataAccessSelfServiceRule(parentId string, input *DataAccessSelfServiceRule) (*DataAccessSelfServiceRule, error) {
	output := DataAccessSelfServiceRule{}
	params := make(map[string]string, 1)
	params["parentId"] = parentId
	return &output, c.postJsonWithParams(DataAccessSelfServiceApiPrefix, &params, input, &output)
}

func (c *Client) UpdateDataAccessSelfServiceRule(id string, input *DataAccessSelfServiceRule) (*DataAccessSelfServiceRule, error) {
	output := DataAccessSelfServiceRule{}
	return &output, c.putJson(DataAccessSelfServiceApiPrefix, "", id, input, &output)
}

func (c *Client) UpdateDataAccessSelfServiceSuspendedStatus(id string, suspend bool) (*DataAccessSelfServiceRule, error) {
	output := DataAccessSelfServiceRule{}
	params := map[string]string{"shouldSuspend": strconv.FormatBool(suspend)}
	return &output, c.putWithParams(DataAccessSelfServiceApiPrefix, "suspend", id, &params, &output)
}

func (c *Client) GetDataAccessSelfServiceRule(id string) (*DataAccessSelfServiceRule, error, int) {
	output := DataAccessSelfServiceRule{}
	err, statusCode := c.getJsonById(DataAccessSelfServiceApiPrefix, "", id, &output)
	return &output, err, statusCode
}

func (c *Client) DeleteDataAccessSelfServiceRule(id string) error {
	return c.delete(DataAccessSelfServiceApiPrefix, id)
}
