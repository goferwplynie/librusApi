package librus

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type User struct {
	ID                       int    `json:"Id"`
	AccountID                string `json:"AccountId"`
	AccountNumericIdentifier int    `json:"AccountNumericIdentifier"`
	FirstName                string `json:"FirstName"`
	LastName                 string `json:"LastName"`
	IsSchoolAdministrator    bool   `json:"IsSchoolAdministrator"`
	IsEmployee               bool   `json:"IsEmployee"`
	GroupID                  int    `json:"GroupId"`
}

// User returns the user (e.g. a teacher) with the given id.
func (c *Client) User(id int) (*User, error) {
	strId := strconv.Itoa(id)

	resp, err := c.Get("Users/" + strId)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userResp struct {
		User User `json:"User"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json:\n%w", err)
	}

	return &userResp.User, nil
}

// Users returns all users visible to the logged in account.
func (c *Client) Users() ([]User, error) {
	resp, err := c.Get("Users")
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
