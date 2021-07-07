package api

type DataStoreLocation struct {
	DataStoreId string `json:"dataStoreId"`
	Location    string `json:"location"`
}

type DataSet struct {
	Name             string              `json:"name"`
	Description      string              `json:"description"`
	OwnersIds        []string            `json:"ownersIds"`
	IncludeLocations []DataStoreLocation `json:"includeLocations"`
	ExcludeLocations []DataStoreLocation `json:"excludeLocations"`
}

type DataSetOutput struct {
	Id               string              `json:"id"`
	Name             string              `json:"name"`
	Description      string              `json:"description"`
	OwnersIds        []string            `json:"ownersIds"`
	IncludeLocations []DataStoreLocation `json:"includeLocations"`
	ExcludeLocations []DataStoreLocation `json:"excludeLocations"`
	DataPolicyId     string              `json:"dataPolicyId"`
}

// GetDataSets - Returns list of data sets
func (c *Client) GetDataSets() (*[]DataSetOutput, error) {
	var output struct {
		Count   int             `json:"count"`
		Records []DataSetOutput `json:"records"`
	}
	return &output.Records, c.getJsonForAccount("/api/dataset", &output)
}

func (c *Client) CreateDataSet(input *DataSet) (*DataSetOutput, error) {
	output := DataSetOutput{}
	return &output, c.postJsonForAccount("/api/dataset", input, &output)
}

func (c *Client) UpdateDataSet(id string, input *DataSet) (*DataSetOutput, error) {
	output := DataSetOutput{}
	return &output, c.putJson("/api/dataset", id, input, &output)
}

// GetDataSet - Returns data set for the given ID
func (c *Client) GetDataSet(id string) (*DataSetOutput, error) {
	output := DataSetOutput{}
	return &output, c.getJson("/api/dataset", id, &output)
}

func (c *Client) DeleteDataSet(id string) error {
	return c.delete("/api/dataset", id)
}
