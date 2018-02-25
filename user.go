package travisci

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

func (c *Client) CurrentUser() (*User, error) {
	url, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}

	url.Path = "/user"
	fmt.Printf("%s\n", url.String())

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Travis-API-Version", "3")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New(response.Status)
	}

	var user User
	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
