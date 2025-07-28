package librusApi

import (
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

type GradeParams struct {
	client                *Client
	SubjectId             int
	StudentId             int
	CategoryId            int
	Semester              int
	AddedById             int
	Grade                 string
	Date                  time.Time
	AddDate               time.Time
	IsConstituent         bool
	IsSemesterGrade       bool
	IsSemesterProposition bool
	IsFinalGrade          bool
	IsFinalProposition    bool
}

func (g Grade) GetTeacher() (*User, error) {
	return GetUser(g.client, g.AddedById)
}
