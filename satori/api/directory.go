package api

const DirectoryGroupApiPrefix = "/api/directory/group"

type DirectoryGroup struct {
	Id          *string                `json:"id,omitempty"`
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	Members     []DirectoryGroupMember `json:"members"`
}

type DirectoryGroupMember struct {
	Type     string  `json:"type"`
	Name     *string `json:"name,omitempty"`
	DsType   *string `json:"dsType,omitempty"`
	Provider *string `json:"provider,omitempty"`
	GroupId  *string `json:"id,omitempty"`
}

func (c *Client) CreateDirectoryGroup(input *DirectoryGroup) (*DirectoryGroup, error) {
	output := DirectoryGroup{}
	return &output, c.postJsonForAccount(DirectoryGroupApiPrefix, input, &output)
}

func (c *Client) UpdateDirectoryGroup(id string, input *DirectoryGroup) (*DirectoryGroup, error) {
	output := DirectoryGroup{}
	return &output, c.putJson(DirectoryGroupApiPrefix, id, input, &output)
}

func (c *Client) GetDirectoryGroup(id string) (*DirectoryGroup, error, int) {
	output := DirectoryGroup{}
	err, statusCode := c.getJsonById(DirectoryGroupApiPrefix, id, &output)
	return &output, err, statusCode
}

func (c *Client) DeleteDirectoryGroup(id string) error {
	return c.delete(DirectoryGroupApiPrefix, id)
}
