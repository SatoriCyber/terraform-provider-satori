package api

const DataStoreApiPrefix = "/api/v1/datastore"

type DataStore struct {
	Name                   string                  `json:"name"`
	Hostname               string                  `json:"hostname"`
	SatoriHostname         string                  `json:"satoriHostname,omitempty"`
	Port                   int                     `json:"port"`
	Type                   string                  `json:"type"`
	DataAccessControllerId string                  `json:"dataAccessControllerId"`
	BaselineSecurityPolicy *BaselineSecurityPolicy `json:"baselineSecurityPolicy,omitempty"`
	IdentityProviderId     string                  `json:"identityProviderId,omitempty"`
	ProjectIds             []string                `json:"projectIds,omitempty"`
	port                   int                     `json:"port"`
	CustomIngressPort      int                     `json:"customIngressPort"`
}

type DataStoreOutput struct {
	Id                     string                  `json:"id"`
	Name                   string                  `json:"name"`
	Hostname               string                  `json:"hostname"`
	SatoriHostname         string                  `json:"satoriHostname,omitempty"`
	Port                   int                     `json:"originPort"`
	CustomIngressPort      int                     `json:"customIngressPort"`
	Type                   string                  `json:"type"`
	ParentId               string                  `json:"parentId"`
	DataPolicyId           string                  `json:"dataPolicyId"`
	DataAccessControllerId string                  `json:"dataAccessControllerId"`
	BaselineSecurityPolicy *BaselineSecurityPolicy `json:"baselineSecurityPolicy,omitempty"`

	IdentityProviderId string   `json:"identityProviderId"`
	ProjectIds         []string `json:"projectIds,omitempty"`
}

type UnassociatedQueriesCategory struct {
	QueryAction string `json:"queryAction"`
}
type UnsupportedQueriesCategory struct {
	QueryAction string `json:"queryAction"`
}
type ExcludedIdentities struct {
	IdentityType string `json:"identityType"`
	Identity     string `json:"identity"`
}
type ExcludedQueryPatterns struct {
	Pattern string `json:"pattern"`
}
type Exclusions struct {
	ExcludedIdentities    []ExcludedIdentities    `json:"excludedIdentities,omitempty"`
	ExcludedQueryPatterns []ExcludedQueryPatterns `json:"excludedQueryPatterns,omitempty"`
}
type BaselineSecurityPolicy struct {
	Type                        string                      `json:"type,omitempty"`
	UnassociatedQueriesCategory UnassociatedQueriesCategory `json:"unassociatedQueriesCategory,omitempty"`
	UnsupportedQueriesCategory  UnsupportedQueriesCategory  `json:"unsupportedQueriesCategory,omitempty"`
	Exclusions                  Exclusions                  `json:"exclusions,omitempty"`
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
	var output DataStoreOutput
	err, statusCode := c.getJsonById(DataStoreApiPrefix, "", id, &output)
	return &output, err, statusCode
}

func (c *Client) DeleteDataStore(id string) error {
	return c.delete(DataStoreApiPrefix, id)
}
