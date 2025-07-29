package librusApi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func GetUser(c *Client, id int) (*User, error) {
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

func GetUsers(c *Client) ([]User, error) {
	resp, err := c.Get("Users", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users struct {
		Users []User `json:"Users"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json")
	}

	return users.Users, nil
}
