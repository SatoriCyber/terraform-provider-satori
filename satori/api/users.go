package api

const UsersApiPrefix = "/api/v1/users"

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func (c *Client) QueryUsers(search *string) ([]User, error) {
	var output struct {
		Count   int    `json:"count"`
		Records []User `json:"records"`
	}
	return output.Records, c.getJsonForAccount(UsersApiPrefix, search, &output)
}
