package api

const DataAccessSelfServiceApiPrefix = "/api/data-access-self-service"

type DataAccessSelfServiceRule struct {
	Id               *string                        `json:"id,omitempty"`
	ParentId         *string                        `json:"parentId,omitempty"`
	Identity         *DataAccessIdentity            `json:"identity,omitempty"`
	AccessLevel      string                         `json:"accessLevel"`
	TimeLimit        DataAccessSelfServiceTimeLimit `json:"timeLimit"`
	UnusedTimeLimit  DataAccessUnusedTimeLimit      `json:"unusedTimeLimit"`
	SecurityPolicies *[]string                      `json:"securityPolicyIds,omitempty"`
}

type DataAccessSelfServiceTimeLimit struct {
	ShouldExpire bool   `json:"shouldExpire"`
	UnitType     string `json:"unitType"`
	Units        int    `json:"units"`
}

func (c *Client) CreateDataAccessSelfServiceRule(parentId string, input *DataAccessSelfServiceRule) (*DataAccessSelfServiceRule, error) {
	output := DataAccessSelfServiceRule{}
	params := make(map[string]string, 1)
	params["parentId"] = parentId
	return &output, c.postJsonWithParams(DataAccessSelfServiceApiPrefix, &params, input, &output)
}

func (c *Client) UpdateDataAccessSelfServiceRule(id string, input *DataAccessSelfServiceRule) (*DataAccessSelfServiceRule, error) {
	output := DataAccessSelfServiceRule{}
	return &output, c.putJson(DataAccessSelfServiceApiPrefix, id, input, &output)
}

func (c *Client) GetDataAccessSelfServiceRule(id string) (*DataAccessSelfServiceRule, error, int) {
	output := DataAccessSelfServiceRule{}
	err, statusCode := c.getJsonById(DataAccessSelfServiceApiPrefix, id, &output)
	return &output, err, statusCode
}

func (c *Client) DeleteDataAccessSelfServiceRule(id string) error {
	return c.delete(DataAccessSelfServiceApiPrefix, id)
}
