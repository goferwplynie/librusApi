package librusApi

import (
	"encoding/json"
	"fmt"
)

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
