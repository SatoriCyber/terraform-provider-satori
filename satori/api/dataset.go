package api

const DataSetApiPrefix = "/api/v1/dataset"

type DataSetLocation struct {
	DataStoreId string                  `json:"dataStoreId"`
	Location    *DataSetGenericLocation `json:"location,omitempty"`
}

type DataSetGenericLocation struct {
	Type   string  `json:"type"`
	Db     *string `json:"db,omitempty"`
	Schema *string `json:"schema,omitempty"`
	Table  *string `json:"table,omitempty"`
}

type DataSet struct {
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	OwnersIds        []string          `json:"ownersIds"`
	IncludeLocations []DataSetLocation `json:"includeLocations"`
	ExcludeLocations []DataSetLocation `json:"excludeLocations"`
}

type DataSetOutput struct {
	Id               string            `json:"id"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	OwnersIds        []string          `json:"ownersIds"`
	IncludeLocations []DataSetLocation `json:"includeLocations"`
	ExcludeLocations []DataSetLocation `json:"excludeLocations"`
	DataPolicyId     string            `json:"dataPolicyId"`
}

func (c *Client) CreateDataSet(input *DataSet) (*DataSetOutput, error) {
	output := DataSetOutput{}
	return &output, c.postJsonForAccount(DataSetApiPrefix, input, &output)
}

func (c *Client) UpdateDataSet(id string, input *DataSet) (*DataSetOutput, error) {
	output := DataSetOutput{}
	return &output, c.putJson(DataSetApiPrefix, "", id, input, &output)
}

func (c *Client) GetDataSet(id string) (*DataSetOutput, error, int) {
	output := DataSetOutput{}
	err, statusCode := c.getJsonById(DataSetApiPrefix, "", id, &output)
	return &output, err, statusCode
}

func (c *Client) DeleteDataSet(id string) error {
	return c.delete(DataSetApiPrefix, id)
}
