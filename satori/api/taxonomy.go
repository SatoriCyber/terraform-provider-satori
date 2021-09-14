package api

const TaxonomyApiPrefix = "/api/v1/taxonomy/custom"
const TaxonomyCategoryApiPrefix = TaxonomyApiPrefix + "/category"
const TaxonomyClassifierApiPrefix = TaxonomyApiPrefix + "/classifier"

type TaxonomyCategory struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	ParentNode  *string `json:"parentNode"`
	Color       string  `json:"color"`
}

type TaxonomyClassifierScope struct {
	DatasetIds       []string          `json:"datasetIds"`
	IncludeLocations []DataSetLocation `json:"includeLocations"`
}

type TaxonomyClassifierValues struct {
	Values          *[]string `json:"values,omitempty"`
	CaseInsensitive bool      `json:"caseInsensitive"`
	Regex           bool      `json:"regex"`
}

type TaxonomyClassifierConfig struct {
	Type                   string                    `json:"type"`
	FieldNamePattern       *string                   `json:"fieldNamePattern,omitempty"`
	SatoriBaseClassifierId *string                   `json:"satoriBaseClassifierId,omitempty"`
	FieldType              *string                   `json:"fieldType,omitempty"`
	Values                 *TaxonomyClassifierValues `json:"values,omitempty"`
	AdditionalCategories   []string                  `json:"additionalSatoriCategoriesToTag"`
}

type TaxonomyClassifier struct {
	Name        string                   `json:"name"`
	Description *string                  `json:"description"`
	ParentNode  string                   `json:"parentNode"`
	Scope       TaxonomyClassifierScope  `json:"scope"`
	Config      TaxonomyClassifierConfig `json:"config"`
}

type TaxonomyCategoryOutput struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Tag         string  `json:"tag"`
	Description *string `json:"description,omitempty"`
	ParentNode  *string `json:"parentNode,omitempty"`
	Color       string  `json:"color"`
}

type TaxonomyClassifierOutput struct {
	Id          string                   `json:"id"`
	Name        string                   `json:"name"`
	Tag         string                   `json:"tag"`
	Description *string                  `json:"description"`
	ParentNode  string                   `json:"parentNode"`
	Scope       TaxonomyClassifierScope  `json:"scope"`
	Config      TaxonomyClassifierConfig `json:"config"`
}

func (c *Client) CreateTaxonomyCategory(input *TaxonomyCategory) (*TaxonomyCategoryOutput, error) {
	output := TaxonomyCategoryOutput{}
	return &output, c.postJsonForAccount(TaxonomyCategoryApiPrefix, input, &output)
}

func (c *Client) CreateTaxonomyClassifier(input *TaxonomyClassifier) (*TaxonomyClassifierOutput, error) {
	output := TaxonomyClassifierOutput{}
	return &output, c.postJsonForAccount(TaxonomyClassifierApiPrefix, input, &output)
}

func (c *Client) UpdateTaxonomyCategory(id string, input *TaxonomyCategory) (*TaxonomyCategoryOutput, error) {
	output := TaxonomyCategoryOutput{}
	return &output, c.putJson(TaxonomyCategoryApiPrefix, "", id, input, &output)
}

func (c *Client) UpdateTaxonomyClassifier(id string, input *TaxonomyClassifier) (*TaxonomyClassifierOutput, error) {
	output := TaxonomyClassifierOutput{}
	return &output, c.putJson(TaxonomyClassifierApiPrefix, "", id, input, &output)
}

func (c *Client) GetTaxonomyCategory(id string) (*TaxonomyCategoryOutput, error, int) {
	output := TaxonomyCategoryOutput{}
	err, statusCode := c.getJsonById(TaxonomyApiPrefix, "", id, &output)
	return &output, err, statusCode
}

func (c *Client) GetTaxonomyClassifier(id string) (*TaxonomyClassifierOutput, error, int) {
	output := TaxonomyClassifierOutput{}
	err, statusCode := c.getJsonById(TaxonomyApiPrefix, "", id, &output)
	return &output, err, statusCode
}

func (c *Client) DeleteTaxonomyNode(id string) error {
	return c.delete(TaxonomyApiPrefix, id)
}
