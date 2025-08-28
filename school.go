package librusApi

import (
	"encoding/json"
	"fmt"
)

// Get school information
func GetSchool(c *Client) (*School, error) {
	resp, err := c.Get("Schools", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var schoolsResp schoolsResponse
	if err := json.NewDecoder(resp.Body).Decode(&schoolsResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %v", err)
	}

	return &schoolsResp.School, nil
}
