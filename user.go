package travisci

func (c *Client) CurrentUser() (*User, error) {
	var user User
	err := c.getJSON("user", &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
