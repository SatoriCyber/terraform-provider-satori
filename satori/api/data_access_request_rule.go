package api

import "strconv"

const DataAccessRequestApiPrefix = "/api/v1/data-access-rule/access-request"

type DataAccessRequestRule struct {
	Id               *string                                  `json:"id,omitempty"`
	ParentId         *string                                  `json:"parentId,omitempty"`
	Suspended        bool                                     `json:"suspended,omitempty"`
	Identity         *DataAccessIdentity                      `json:"identity,omitempty"`
	AccessLevel      string                                   `json:"accessLevel"`
	TimeLimit        DataAccessSelfServiceAndRequestTimeLimit `json:"timeLimit"`
	UnusedTimeLimit  DataAccessUnusedTimeLimit                `json:"unusedTimeLimit"`
	SecurityPolicies *[]string                                `json:"securityPolicyIds,omitempty"`
}

func (c *Client) CreateDataAccessRequestRule(parentId string, input *DataAccessRequestRule) (*DataAccessRequestRule, error) {
	output := DataAccessRequestRule{}
	params := make(map[string]string, 1)
	params["parentId"] = parentId
	return &output, c.postJsonWithParams(DataAccessRequestApiPrefix, &params, input, &output)
}

func (c *Client) UpdateDataAccessRequestRule(id string, input *DataAccessRequestRule) (*DataAccessRequestRule, error) {
	output := DataAccessRequestRule{}
	return &output, c.putJson(DataAccessRequestApiPrefix, "", id, input, &output)
}

func (c *Client) UpdateDataAccessRequestSuspendedStatus(id string, suspend bool) (*DataAccessRequestRule, error) {
	output := DataAccessRequestRule{}
	params := map[string]string{"shouldSuspend": strconv.FormatBool(suspend)}
	return &output, c.putWithParams(DataAccessRequestApiPrefix, "suspend", id, &params, &output)
}

func (c *Client) GetDataAccessRequestRule(id string) (*DataAccessRequestRule, error, int) {
	output := DataAccessRequestRule{}
	err, statusCode := c.getJsonById(DataAccessRequestApiPrefix, "", id, &output)
	return &output, err, statusCode
}

func (c *Client) DeleteDataAccessRequestRule(id string) error {
	return c.delete(DataAccessRequestApiPrefix, id)
}
