package librusApi

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
