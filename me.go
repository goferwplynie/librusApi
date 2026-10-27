package librus

import "context"

type meResponse struct {
	Me UserInfo `json:"Me"`
}

type UserInfo struct {
	Account struct {
		ID                 int    `json:"Id"`
		UserID             int    `json:"UserId"`
		FirstName          string `json:"FirstName"`
		LastName           string `json:"LastName"`
		Email              string `json:"Email"`
		GroupID            int    `json:"GroupId"`
		IsActive           bool   `json:"IsActive"`
		Login              string `json:"Login"`
		IsPremium          bool   `json:"IsPremium"`
		IsPremiumDemo      bool   `json:"IsPremiumDemo"`
		ExpiredPremiumDate int    `json:"ExpiredPremiumDate"`
	} `json:"Account"`
	Refresh int `json:"Refresh"`
	User    struct {
		FirstName string `json:"FirstName"`
		LastName  string `json:"LastName"`
	} `json:"User"`
	Class struct {
		ID int `json:"Id"`
	} `json:"Class"`
}

// Me returns information about the logged in account.
func (c *Client) Me(ctx context.Context) (*UserInfo, error) {
	var resp meResponse
	if err := c.getJSON(ctx, "Me", &resp); err != nil {
		return nil, err
	}
	return &resp.Me, nil
}
