package librusApi

import (
	"encoding/json"
	"fmt"
)

// get current user information
func GetUserInfo(c *Client) (*UserInfo, error) {
	resp, err := c.Get("Me", nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var responseObject meResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseObject); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json:\n%w", err)
	}
	return &responseObject.Me, nil
}
