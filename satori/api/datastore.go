package api

const DataStoreApiPrefix = "/api/v1/datastore"

type DataStoreLocation struct {
	DataStoreId string                  `json:"dataStoreId"`
	Location    *DataSetGenericLocation `json:"location,omitempty"`
}

type DataStoreGenericLocation struct {
	Type   string  `json:"type"`
	Db     *string `json:"db,omitempty"`
	Schema *string `json:"schema,omitempty"`
	Table  *string `json:"table,omitempty"`
}

type DataStore struct {
	Name             string              `json:"name"`
	Description      string              `json:"description"`
	OwnersIds        []string            `json:"ownersIds"`
	IncludeLocations []DataStoreLocation `json:"includeLocations"`
	ExcludeLocations []DataStoreLocation `json:"excludeLocations"`
}

type DataStoreOutput struct {
	Id               string              `json:"id"`
	Name             string              `json:"name"`
	Description      string              `json:"description"`
	OwnersIds        []string            `json:"ownersIds"`
	IncludeLocations []DataStoreLocation `json:"includeLocations"`
	ExcludeLocations []DataStoreLocation `json:"excludeLocations"`
	DataPolicyId     string              `json:"dataPolicyId"`
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
