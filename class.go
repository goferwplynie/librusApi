package librus

import "context"

type classResponse struct {
	Class Class  `json:"Class"`
	URL   string `json:"Url"`
}

type Class struct {
	EndFirstSemester string       `json:"EndFirstSemester"`
	ClassTutors      []ClassTutor `json:"ClassTutors"`
	Number           int64        `json:"Number"`
	ClassTutor       ClassTutor   `json:"ClassTutor"`
	Symbol           string       `json:"Symbol"`
	EndSchoolYear    string       `json:"EndSchoolYear"`
	ID               int64        `json:"Id"`
	Unit             ClassTutor   `json:"Unit"`
	BeginSchoolYear  string       `json:"BeginSchoolYear"`
}

type ClassTutor struct {
	ID  int64  `json:"Id"`
	URL string `json:"Url"`
}

// Class returns information about the logged in student's class.
func (c *Client) Class(ctx context.Context) (*Class, error) {
	var resp classResponse
	if err := c.getJSON(ctx, "Classes", &resp); err != nil {
		return nil, err
	}
	return &resp.Class, nil
}
