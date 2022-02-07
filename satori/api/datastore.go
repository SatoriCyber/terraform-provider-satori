package api

const DataStoreApiPrefix = "/api/v1/datastore"

type DataStore struct {
	Name                   string                 `json:"name"`
	Hostname               string                 `json:"hostname"`
	Port                   int                    `json:"port"`
	Type                   string                 `json:"type"`
	Rules                  []map[string]string    `json:"rules"`
	DataAccessControllerId string                 `json:"dataAccessControllerId"`
	BaselineSecurityPolicy BaselineSecurityPolicy `json:"baselineSecurityPolicy"`
	IdentityProviderId     string                 `json:"identityProviderId"`
	ProjectIds             []string               `json:"projectIds,omitempty"`
	port                   int                    `json:"port"`
	CustomIngressPort      int                    `json:"customIngressPort"`
}

type DataStoreOutput struct {
	Id                     string                 `json:"Id"`
	Name                   string                 `json:"Name"`
	Hostname               string                 `json:"Hostname"`
	Port                   int                    `json:"originPort"`
	CustomIngressPort      int                    `json:"customIngressPort"`
	Type                   string                 `json:"type"`
	ParentId               string                 `json:"parentId"`
	DataPolicyId           string                 `json:"dataPolicyId"`
	DataAccessControllerId string                 `json:"dataAccessControllerId"`
	BaselineSecurityPolicy BaselineSecurityPolicy `json:"baselineSecurityPolicy"`

	IdentityProviderId string   `json:"identityProviderId"`
	ProjectIds         []string `json:"projectIds"`
}

type vUnassociatedQueriesCategory struct {
}
type vUnsupportedQueriesCategory struct {
}
type vExclusions struct {
}
type BaselineSecurityPolicy struct {
	Type                        string `json:"type"`
	unassociatedQueriesCategory vUnassociatedQueriesCategory
	unsupportedQueriesCategory  vUnsupportedQueriesCategory
	exclusions                  vExclusions
}

func (c *Client) CreateDataStore(input *DataStore) (*DataStoreOutput, error) {
	output := DataStoreOutput{}
	return &output, c.postJsonForAccount(DataStoreApiPrefix, input, &output)
}

func (c *Client) UpdateDataStore(id string, input *DataStore) (*DataStoreOutput, error) {
	output := DataStoreOutput{}
	return &output, c.putJson(DataStoreApiPrefix, "", id, input, &output)
}

func (c *Client) GetDataStore(id string) (*DataStoreOutput, error, int) {
	output := DataStoreOutput{}
	err, statusCode := c.getJsonById(DataStoreApiPrefix, "", id, &output)
	return &output, err, statusCode
}

func (c *Client) DeleteDataStore(id string) error {
	return c.delete(DataStoreApiPrefix, id)
}
