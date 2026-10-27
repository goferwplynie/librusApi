package librus

import (
	"context"
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
func (c *Client) User(ctx context.Context, id int) (*User, error) {
	var resp struct {
		User User `json:"User"`
	}
	if err := c.getJSON(ctx, "Users/"+strconv.Itoa(id), &resp); err != nil {
		return nil, err
	}
	return &resp.User, nil
}

// Users returns all users visible to the logged in account.
func (c *Client) Users(ctx context.Context) ([]User, error) {
	var resp struct {
		Users []User `json:"Users"`
	}
	if err := c.getJSON(ctx, "Users", &resp); err != nil {
		return nil, err
	}
	return resp.Users, nil
}
