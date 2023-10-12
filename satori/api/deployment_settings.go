package api

import (
	"fmt"
)

const DataAccessControllerApiPrefix = "/api/data-access-controllers"
const GoogleServiceAccountApiSuffix = "/google-sa"

type DeploymentSettings struct {
	DacId string  `json:"dacId"`
	GSA   *string `json:"gsa"`
}

type GSAResponse struct {
	GSA *string `json:"gsa"`
}

func (c *Client) QueryDeploymentSettings(dacId *string) (*DeploymentSettings, error) {
	var output GSAResponse

	var deploymentSettings DeploymentSettings

	params := map[string]string{"id": *dacId}
	err := c.getJsonForAccountWithParams(DataAccessControllerApiPrefix+GoogleServiceAccountApiSuffix, &params, &output)

	if err != nil {
		fmt.Println("Failed to get DAC's SA")
		return nil, err
	}

	fmt.Println("This is the SA output")
	fmt.Println(output)

	deploymentSettings.DacId = *dacId
	deploymentSettings.GSA = output.GSA

	return &deploymentSettings, nil
}
