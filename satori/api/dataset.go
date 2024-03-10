package api

const DataSetApiPrefix = "/api/v1/dataset"

type DataSetLocation struct {
	DataStoreId string                  `json:"dataStoreId"`
	Location    *DataSetGenericLocation `json:"location,omitempty"`
}

// DataSetGenericLocation
// Location types with their parts:
// Relational - Db, Schema, Table
// MySql - Db, Table
// Athena - Catalog, Db, Table
// Mongo - Db, Collection
// S3 - Bucket, ObjectKey
// /**
type DataSetGenericLocation struct {
	Type       string  `json:"type"`
	Db         *string `json:"db,omitempty"`
	Schema     *string `json:"schema,omitempty"`
	Table      *string `json:"table,omitempty"`
	Catalog    *string `json:"catalog,omitempty"`
	Collection *string `json:"collection,omitempty"`
	Bucket     *string `json:"bucket,omitempty"`
	ObjectKey  *string `json:"objectKey,omitempty"`
}

type DataSet struct {
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	OwnersIds        []string          `json:"ownersIds"`
	IncludeLocations []DataSetLocation `json:"includeLocations"`
	ExcludeLocations []DataSetLocation `json:"excludeLocations"`
	// data policy
	PermissionsEnabled bool               `json:"permissionsEnabled"`
	CustomPolicy       CustomPolicy       `json:"customPolicy"`
	SecurityPolicies   SecurityPolicies   `json:"defaultSecurityPolicies"`
	Approvers          []ApproverIdentity `json:"approvers"`
}

type DataSetOutput struct {
	DataSet
	Id           string `json:"id"`
	DataPolicyId string `json:"dataPolicyId"`
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
