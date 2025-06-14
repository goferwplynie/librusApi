package users

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/goferwplynie/librusApi/client"
)

func GetUser(c *client.Client, id int) (*User, error) {
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
