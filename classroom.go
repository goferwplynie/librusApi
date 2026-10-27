package librus

import "context"

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

// Classrooms returns all classrooms listed in the school.
func (c *Client) Classrooms(ctx context.Context) ([]Classroom, error) {
	var resp classroomsResponse
	if err := c.getJSON(ctx, "Classrooms", &resp); err != nil {
		return nil, err
	}
	return resp.Classrooms, nil
}
