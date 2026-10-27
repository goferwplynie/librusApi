package librus

import (
	"context"
	"time"
)

type Grade struct {
	client                *Client
	ID                    int
	LessonID              int
	CommentIDs            []int
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

// Category returns the category of the grade. To resolve many grades at once
// fetch Client.GradeCategories and match by CategoryID instead.
func (g Grade) Category(ctx context.Context) (*GradeCategory, error) {
	return g.client.GradeCategory(ctx, g.CategoryID)
}

// Comments returns the teacher's comments to the grade, one request per comment.
// Grades without comments return nil without any request.
func (g Grade) Comments(ctx context.Context) ([]GradeComment, error) {
	var comments []GradeComment
	for _, id := range g.CommentIDs {
		comment, err := g.client.GradeComment(ctx, id)
		if err != nil {
			return nil, err
		}
		comments = append(comments, *comment)
	}
	return comments, nil
}

// Subject returns the subject the grade was given for.
func (g Grade) Subject(ctx context.Context) (*Subject, error) {
	return g.client.Subject(ctx, g.SubjectID)
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
		Comments              []ref  `json:"Comments,omitempty"`
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
		date, err := parseDate(v.Date)
		if err != nil {
			return nil, err
		}
		addDate, err := parseDateTime(v.AddDate)
		if err != nil {
			return nil, err
		}

		var commentIDs []int
		for _, comment := range v.Comments {
			commentIDs = append(commentIDs, int(comment.ID))
		}

		grades = append(grades, &Grade{
			client:                c,
			ID:                    v.ID,
			LessonID:              v.Lesson.ID,
			CommentIDs:            commentIDs,
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
