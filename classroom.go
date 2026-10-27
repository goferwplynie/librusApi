package librusApi

import (
	"encoding/json"
	"fmt"
)

type classroomsResponse struct {
	Classrooms []Classroom `json:"Classrooms"`
	URL        string      `json:"Url"`
}

type Classroom struct {
	SchoolCommonRoom bool    `json:"SchoolCommonRoom"`
	Description      *string `json:"Description,omitempty"`
	Symbol           string  `json:"Symbol"`
	Size             int64   `json:"Size"`
	ID               int64   `json:"Id"`
	Name             string  `json:"Name"`
}

// Get all listed classrooms
func GetClassrooms(c *Client) ([]Classroom, error) {
	resp, err := c.Get("Classrooms", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var classroomsResp classroomsResponse
	if err := json.NewDecoder(resp.Body).Decode(&classroomsResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %v", err)
	}

	return classroomsResp.Classrooms, nil
}
