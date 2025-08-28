package librusApi

import (
	"encoding/json"
	"fmt"
)

// Get class information
func GetClass(c *Client) (*Class, error) {
	resp, err := c.Get("Classes", nil)
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
