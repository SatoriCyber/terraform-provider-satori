package api

const DataStoreApiPrefix = "/api/v1/datastore"

type DataStore struct {
	Name                   string              `json:"name"`
	Hostname               string              `json:"hostname"`
	Port                   int                 `json:"port"`
	Type                   string              `json:"type"`
	Rules                  []map[string]string `json:"rules"`
	DataAccessControllerId string              `json:"dataAccessControllerId"`
	CustomIngressPort      int                 `json:"customIngressPort"`
	BaselineSecurityPolicy []map[string]string `json:"baselineSecurityPolicy"`
	IdentityProviderId     string              `json:"identityProviderId"`
	ProjectIds             []map[string]string `json:"projectIds"`
}

type DataStoreOutput struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	OwnersIds   []string `json:"ownersIds"`

	DataPolicyId string `json:"dataPolicyId"`
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
