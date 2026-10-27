package librus

import (
	"context"
	"strconv"
)

type Subject struct {
	ID                int    `json:"Id"`
	Name              string `json:"Name"`
	Short             string `json:"Short"`
	No                int    `json:"No"`
	IsExtracurricular bool   `json:"IsExtracurricular"`
	IsBlockLesson     bool   `json:"IsBlockLesson"`
}

// Subject returns the subject with the given id.
func (c *Client) Subject(ctx context.Context, id int) (*Subject, error) {
	var resp struct {
		Subject Subject `json:"Subject"`
	}
	if err := c.getJSON(ctx, "Subjects/"+strconv.Itoa(id), &resp); err != nil {
		return nil, err
	}
	return &resp.Subject, nil
}

// Subjects returns all subjects of the school.
func (c *Client) Subjects(ctx context.Context) ([]Subject, error) {
	var resp struct {
		Subjects []Subject `json:"Subjects"`
	}
	if err := c.getJSON(ctx, "Subjects", &resp); err != nil {
		return nil, err
	}
	return resp.Subjects, nil
}
