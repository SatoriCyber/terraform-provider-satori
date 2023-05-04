package api

const UsersV1ApiPrefix = "/api/v1/users"

const UsersApiPrefix = "/api/users"
const UserProfileSuffix = "/profile"
const UserAttributesApiSuffix = "/attributes/custom"

type User struct {
	Id         string            `json:"id"`
	Email      string            `json:"email"`
	Attributes map[string]string `json:"customAttributes,omitempty"`
}

type UserWithSatoriAttributes struct {
	Id               string                 `json:"id"`
	CustomAttributes map[string]interface{} `json:"customAttributes"`
}

type SatoriAttributes struct {
	UserId     string                 `json:"userId"`
	Attributes map[string]interface{} `json:"attributes"`
}

func (c *Client) QueryUsers(search *string) ([]User, error) {
	var output struct {
		Count   int    `json:"count"`
		Records []User `json:"records"`
	}
	err, _ := c.getJsonForAccount(UsersV1ApiPrefix, search, &output)
	return output.Records, err
}

type UserProfileAttributes struct {
	Id   string                 `json:"id"`
	Attr map[string]interface{} `json:"customAttributes"`
}

func (c *Client) GetUserProfile(userId *string) (UserProfileAttributes, error, int) {
	output := UserProfileAttributes{}
	profilePath := UsersApiPrefix + "/" + *userId + UserProfileSuffix
	err, responseStatus := c.getJsonForAccount(profilePath, nil, &output)
	return output, err, responseStatus
}

func (c *Client) CreateUserCustomAttributes(input *SatoriAttributes) (*SatoriAttributes, error) {
	output := SatoriAttributes{}
	userResponse := UserWithSatoriAttributes{}
	errResponse := c.putJson(UsersApiPrefix, UserAttributesApiSuffix, input.UserId, input.Attributes, &userResponse)
	return convertResponseToOutput(&output, &userResponse), errResponse
}

func (c *Client) UpdateUserCustomAttributes(input *SatoriAttributes) (*SatoriAttributes, error) {
	output := SatoriAttributes{}
	userResponse := UserWithSatoriAttributes{}
	errResponse := c.putJson(UsersApiPrefix, UserAttributesApiSuffix, input.UserId, input.Attributes, &userResponse)
	return convertResponseToOutput(&output, &userResponse), errResponse
}

func (c *Client) DeleteUserCustomAttributes(input *SatoriAttributes) (*SatoriAttributes, error) {
	output := SatoriAttributes{}
	userResponse := UserWithSatoriAttributes{}
	// Updates the attributes to be empty body '{}'
	errResponse := c.putJson(UsersApiPrefix, UserAttributesApiSuffix, input.UserId, make(map[string]string), &userResponse)

	return convertResponseToOutput(&output, &userResponse), errResponse
}

// Function to convert from the response of the update satori custom attributes API
// to the proper SatoriAttributes terraform schema
func convertResponseToOutput(output *SatoriAttributes, response *UserWithSatoriAttributes) *SatoriAttributes {
	output.UserId = response.Id
	output.Attributes = response.CustomAttributes
	return output
}
