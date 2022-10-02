package api

const DataAccessControllerPrefix = "/api/v1/data-access-controllers"

type DataAccessController struct {
	Id            string `json:"id"`
	Type          string `json:"type"`
	Region        string `json:"region"`
	CloudProvider string `json:"cloudProvider"`
	UniqueName    string `json:"uniqueName"`
}

func (c *Client) QueryDataAccessControllers(dacType *string, region *string, cloudProvider *string, uniqueName *string) ([]DataAccessController, error) {
	var output struct {
		Count   int                    `json:"count"`
		Records []DataAccessController `json:"records"`
	}
	params := map[string]string{"type": *dacType, "region": *region, "cloudProvider": *cloudProvider, "uniqueName": *uniqueName}
	return output.Records, c.getJsonForAccountWithParams(DataAccessControllerPrefix, &params, &output)
}
