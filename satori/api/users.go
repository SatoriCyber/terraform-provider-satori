package api

const UsersV1ApiPrefix = "/api/v1/users"

const UsersApiPrefix = "/api/users"
const UserProfileSuffix = "/profile"
const UserAttributesApiSuffix = "/attributes/custom"

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type UserWithCustomAttributes struct {
	UserId           string                 `json:"id"`
	CustomAttributes map[string]interface{} `json:"customAttributes"`
}

func (c *Client) QueryUsers(search *string) ([]User, error) {
	var output struct {
		Count   int    `json:"count"`
		Records []User `json:"records"`
	}
	err, _ := c.getJsonForAccount(UsersV1ApiPrefix, search, &output)
	return output.Records, err
}

func (c *Client) GetUserProfile(userId *string) (*UserWithCustomAttributes, error, int) {
	output := UserWithCustomAttributes{}
	profilePath := UsersApiPrefix + "/" + *userId + UserProfileSuffix
	err, responseStatus := c.getJsonForAccount(profilePath, nil, &output)
	return &output, err, responseStatus
}

func (c *Client) CreateUserCustomAttributes(input *UserWithCustomAttributes) (*UserWithCustomAttributes, error) {
	output := UserWithCustomAttributes{}
	errResponse := c.putJson(UsersApiPrefix, UserAttributesApiSuffix, input.UserId, input.CustomAttributes, &output)
	return &output, errResponse
}

func (c *Client) UpdateUserCustomAttributes(input *UserWithCustomAttributes) (*UserWithCustomAttributes, error) {
	output := UserWithCustomAttributes{}
	errResponse := c.putJson(UsersApiPrefix, UserAttributesApiSuffix, input.UserId, input.CustomAttributes, &output)
	return &output, errResponse
}

func (c *Client) DeleteUserCustomAttributes(input *UserWithCustomAttributes) (*UserWithCustomAttributes, error) {
	output := UserWithCustomAttributes{}
	// Updates the attributes to be empty body '{}'
	errResponse := c.putJson(UsersApiPrefix, UserAttributesApiSuffix, input.UserId, make(map[string]string), &output)
	return &output, errResponse
}
