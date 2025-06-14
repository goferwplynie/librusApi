package grades

import (
	"time"

	"github.com/goferwplynie/librusApi/client"
	"github.com/goferwplynie/librusApi/users"
)

type Grade struct {
	client                *client.Client
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
	client                *client.Client
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

func (g Grade) GetTeacher() (*users.User, error) {
	return users.GetUser(g.client, g.AddedById)
}
