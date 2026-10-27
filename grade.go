package librus

import (
	"encoding/json"
	"fmt"
	"time"
)

type Grade struct {
	client                *Client
	SubjectId             int
	StudentId             int
	CategoryId            int
	AddedById             int
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
func (g Grade) Teacher() (*User, error) {
	return g.client.User(g.AddedById)
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
func (c *Client) Grades() ([]*Grade, error) {
	grades := make([]*Grade, 0)
	gradesJsonResp, err := c.Get("Grades")
	if err != nil {
		return nil, err
	}
	defer gradesJsonResp.Body.Close()

	var gradesResp gradesResponse
	if err := json.NewDecoder(gradesJsonResp.Body).Decode(&gradesResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json:\n%w", err)
	}

	for _, v := range gradesResp.Grades {
		date, err := time.Parse(time.DateOnly, v.Date)
		if err != nil {
			return nil, fmt.Errorf("failed parsing time:\n%w", err)
		}
		addDate, err := time.Parse(time.DateTime, v.AddDate)
		if err != nil {
			return nil, fmt.Errorf("failed parsing time:\n%w", err)
		}

		grade := &Grade{
			client:                c,
			SubjectId:             v.Subject.ID,
			StudentId:             v.Student.ID,
			CategoryId:            v.Category.ID,
			Semester:              v.Semester,
			AddedById:             v.AddedBy.ID,
			Grade:                 v.Grade,
			Date:                  date,
			AddDate:               addDate,
			IsConstituent:         v.IsConstituent,
			IsSemesterGrade:       v.IsSemester,
			IsSemesterProposition: v.IsSemesterProposition,
			IsFinalGrade:          v.IsFinal,
			IsFinalProposition:    v.IsFinalProposition,
		}

		grades = append(grades, grade)

	}

	return grades, nil
}
