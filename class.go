package librus

import (
	"encoding/json"
	"fmt"
)

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
func (c *Client) Class() (*Class, error) {
	resp, err := c.Get("Classes")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var classResp classResponse
	if err := json.NewDecoder(resp.Body).Decode(&classResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %v", err)
	}

	return &classResp.Class, nil
}
