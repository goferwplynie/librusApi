package librusApi

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
