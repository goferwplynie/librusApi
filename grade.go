package librus

import (
	"context"
	"fmt"
	"time"
)

type Grade struct {
	client                *Client
	SubjectID             int
	StudentID             int
	CategoryID            int
	AddedByID             int
	Grade                 string
	Date                  time.Time
	AddDate               time.Time
	Semester              int
	IsConstituent         bool
	IsSemesterGrade       bool
	IsSemesterProposition bool
	IsFinalGrade          bool
	IsFinalProposition    bool
}

// Teacher returns the user who added the grade.
func (g Grade) Teacher(ctx context.Context) (*User, error) {
	return g.client.User(ctx, g.AddedByID)
}

type gradesResponse struct {
	Grades []struct {
		ID     int `json:"Id"`
		Lesson struct {
			ID  int    `json:"Id"`
			URL string `json:"Url"`
		} `json:"Lesson"`
		Subject struct {
			ID  int    `json:"Id"`
			URL string `json:"Url"`
		} `json:"Subject"`
		Student struct {
			ID  int    `json:"Id"`
			URL string `json:"Url"`
		} `json:"Student"`
		Category struct {
			ID  int    `json:"Id"`
			URL string `json:"Url"`
		} `json:"Category"`
		AddedBy struct {
			ID  int    `json:"Id"`
			URL string `json:"Url"`
		} `json:"AddedBy"`
		Grade                 string `json:"Grade"`
		Date                  string `json:"Date"`
		AddDate               string `json:"AddDate"`
		Semester              int    `json:"Semester"`
		IsConstituent         bool   `json:"IsConstituent"`
		IsSemester            bool   `json:"IsSemester"`
		IsSemesterProposition bool   `json:"IsSemesterProposition"`
		IsFinal               bool   `json:"IsFinal"`
		IsFinalProposition    bool   `json:"IsFinalProposition"`
		Comments              []struct {
			ID  int    `json:"Id"`
			URL string `json:"Url"`
		} `json:"Comments,omitempty"`
	} `json:"Grades"`
}

// Grades returns all grades of the logged in student.
func (c *Client) Grades(ctx context.Context) ([]*Grade, error) {
	var resp gradesResponse
	if err := c.getJSON(ctx, "Grades", &resp); err != nil {
		return nil, err
	}

	grades := make([]*Grade, 0, len(resp.Grades))
	for _, v := range resp.Grades {
		date, err := time.Parse(time.DateOnly, v.Date)
		if err != nil {
			return nil, fmt.Errorf("librus: parsing grade date: %w", err)
		}
		addDate, err := time.Parse(time.DateTime, v.AddDate)
		if err != nil {
			return nil, fmt.Errorf("librus: parsing grade add date: %w", err)
		}

		grades = append(grades, &Grade{
			client:                c,
			SubjectID:             v.Subject.ID,
			StudentID:             v.Student.ID,
			CategoryID:            v.Category.ID,
			Semester:              v.Semester,
			AddedByID:             v.AddedBy.ID,
			Grade:                 v.Grade,
			Date:                  date,
			AddDate:               addDate,
			IsConstituent:         v.IsConstituent,
			IsSemesterGrade:       v.IsSemester,
			IsSemesterProposition: v.IsSemesterProposition,
			IsFinalGrade:          v.IsFinal,
			IsFinalProposition:    v.IsFinalProposition,
		})
	}

	return grades, nil
}
