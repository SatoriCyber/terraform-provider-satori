package api

const DataAccessControllerPrefix = "/api/v1/data-access-controllers"

type DataAccessController struct {
	Id            string   `json:"id"`
	Type          string   `json:"type"`
	Region        string   `json:"region"`
	CloudProvider string   `json:"cloudProvider"`
	Ips           []string `json:"ips"`
}

func (c *Client) QueryDataAccessControllers(dacType *string, region *string, cloudProvider *string) ([]DataAccessController, error) {
	var output struct {
		Count   int                    `json:"count"`
		Records []DataAccessController `json:"records"`
	}
	params := map[string]string{"type": *dacType, "region": *region, "cloudProvider": *cloudProvider}
	return output.Records, c.getJsonForAccountWithParams(DataAccessControllerPrefix, &params, &output)
}

func (c *Client) QueryDataAccessControllerById(id *string) (DataAccessController, error) {
	output := DataAccessController{}
	err, _ := c.getJsonById(DataAccessControllerPrefix, "", *id, &output)
	return output, err
}
