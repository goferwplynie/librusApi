package librus

import (
	"context"
	"strconv"
	"time"
)

// Substitution is a change in the timetable: a lesson that was canceled,
// moved or taken by another teacher or subject. Timetable shows the same
// changes together with regular lessons; this lists only the changes.
type Substitution struct {
	client *Client
	ID     int

	// Date, LessonNo, TeacherID, SubjectID and ClassroomSymbol describe the lesson
	// after the change. They're zero for a canceled lesson.
	Date            time.Time
	LessonNo        int
	TeacherID       int
	SubjectID       int
	ClassroomSymbol string

	// Org* describe the lesson as it was planned.
	OrgDate            time.Time
	OrgLessonNo        int
	OrgTeacherID       int
	OrgSubjectID       int
	OrgClassroomSymbol string

	IsCanceled bool
	IsShifted  bool // the lesson was moved to another time
}

// Teacher returns the teacher of the lesson after the change, or nil if the
// lesson was canceled.
func (s Substitution) Teacher(ctx context.Context) (*User, error) {
	if s.TeacherID == 0 {
		return nil, nil
	}
	return s.client.User(ctx, s.TeacherID)
}

// Subject returns the subject of the lesson after the change, or nil if the
// lesson was canceled.
func (s Substitution) Subject(ctx context.Context) (*Subject, error) {
	if s.SubjectID == 0 {
		return nil, nil
	}
	return s.client.Subject(ctx, s.SubjectID)
}

// OrgTeacher returns the teacher of the planned lesson.
func (s Substitution) OrgTeacher(ctx context.Context) (*User, error) {
	return s.client.User(ctx, s.OrgTeacherID)
}

// OrgSubject returns the subject of the planned lesson.
func (s Substitution) OrgSubject(ctx context.Context) (*Subject, error) {
	return s.client.Subject(ctx, s.OrgSubjectID)
}

type substitutionJSON struct {
	ID                 flexInt `json:"Id"`
	Date               string  `json:"Date"`
	LessonNo           flexInt `json:"LessonNo"`
	Teacher            ref     `json:"Teacher"`
	Subject            ref     `json:"Subject"`
	ClassroomSymbol    *string `json:"ClassroomSymbol"`
	OrgDate            string  `json:"OrgDate"`
	OrgLessonNo        flexInt `json:"OrgLessonNo"`
	OrgTeacher         ref     `json:"OrgTeacher"`
	OrgSubject         ref     `json:"OrgSubject"`
	OrgClassroomSymbol *string `json:"OrgClassroomSymbol"`
	IsCancelled        bool    `json:"IsCancelled"`
	IsShifted          bool    `json:"IsShifted"`
}

func (c *Client) toSubstitution(j substitutionJSON) (*Substitution, error) {
	s := &Substitution{
		client:       c,
		ID:           int(j.ID),
		LessonNo:     int(j.LessonNo),
		TeacherID:    int(j.Teacher.ID),
		SubjectID:    int(j.Subject.ID),
		OrgLessonNo:  int(j.OrgLessonNo),
		OrgTeacherID: int(j.OrgTeacher.ID),
		OrgSubjectID: int(j.OrgSubject.ID),
		IsCanceled:   j.IsCancelled,
		IsShifted:    j.IsShifted,
	}
	if j.ClassroomSymbol != nil {
		s.ClassroomSymbol = *j.ClassroomSymbol
	}
	if j.OrgClassroomSymbol != nil {
		s.OrgClassroomSymbol = *j.OrgClassroomSymbol
	}
	var err error
	if j.Date != "" {
		if s.Date, err = parseDate(j.Date); err != nil {
			return nil, err
		}
	}
	if j.OrgDate != "" {
		if s.OrgDate, err = parseDate(j.OrgDate); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// Substitution returns the timetable change with the given id.
func (c *Client) Substitution(ctx context.Context, id int) (*Substitution, error) {
	var resp struct {
		Substitution substitutionJSON `json:"Substitution"`
	}
	if err := c.getJSON(ctx, "Calendars/Substitutions/"+strconv.Itoa(id), &resp); err != nil {
		return nil, err
	}
	return c.toSubstitution(resp.Substitution)
}

// Substitutions returns the timetable changes (canceled, moved and substituted
// lessons), in the order given by Librus.
func (c *Client) Substitutions(ctx context.Context) ([]*Substitution, error) {
	var resp struct {
		Substitutions []substitutionJSON `json:"Substitutions"`
	}
	if err := c.getJSON(ctx, "Calendars/Substitutions", &resp); err != nil {
		return nil, err
	}
	substitutions := make([]*Substitution, 0, len(resp.Substitutions))
	for _, j := range resp.Substitutions {
		s, err := c.toSubstitution(j)
		if err != nil {
			return nil, err
		}
		substitutions = append(substitutions, s)
	}
	return substitutions, nil
}
