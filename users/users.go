package users

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/goferwplynie/librusApi"
)

func GetUser(c *librusApi.Client, id int) (*User, error) {
	strId := strconv.Itoa(id)

	resp, err := c.Get("Users/"+strId, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json:\n%w", err)
	}

	return &user, nil
}

func GetUsers(c *librusApi.Client) ([]*User, error) {
	resp, err := c.Get("Users", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users struct {
		Users []*User `json:"Users"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json")
	}

	return users.Users, nil
}
