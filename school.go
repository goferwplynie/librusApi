package librus

import "context"

type schoolsResponse struct {
	School School `json:"School"`
	URL    string `json:"Url"`
}

type School struct {
	Email              string         `json:"Email"`
	SurnameHeadTeacher string         `json:"SurnameHeadTeacher"`
	Service            Service        `json:"Service"`
	NameHeadTeacher    string         `json:"NameHeadTeacher"`
	SchoolYear         int64          `json:"SchoolYear"`
	ApartmentNumber    string         `json:"ApartmentNumber"`
	Name               string         `json:"Name"`
	Project            int64          `json:"Project"`
	LessonsRange       []LessonsRange `json:"LessonsRange"`
	Town               string         `json:"Town"`
	State              string         `json:"State"`
	BuildingNumber     string         `json:"BuildingNumber"`
	Street             string         `json:"Street"`
	VocationalSchool   int64          `json:"VocationalSchool"`
	PhoneNumber        string         `json:"PhoneNumber"`
	ID                 int64          `json:"Id"`
	PostCode           string         `json:"PostCode"`
}

type LessonsRange struct {
	RawFrom *int64  `json:"RawFrom,omitempty"`
	RawTo   *int64  `json:"RawTo,omitempty"`
	From    *string `json:"From,omitempty"`
	To      *string `json:"To,omitempty"`
}

type Service struct {
	Variant   Variant `json:"Variant"`
	IsSynergy bool    `json:"IsSynergy"`
	URL       string  `json:"Url"`
}

type Variant struct {
	DisplayName string `json:"DisplayName"`
	Name        string `json:"Name"`
}

// School returns information about the school.
func (c *Client) School(ctx context.Context) (*School, error) {
	var resp schoolsResponse
	if err := c.getJSON(ctx, "Schools", &resp); err != nil {
		return nil, err
	}
	return &resp.School, nil
}
