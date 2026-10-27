package librus

import (
	"context"
	"strconv"
)

// Lesson is a lesson of a class, i.e. a subject taught by a teacher. Grades and
// attendance point at lessons.
type Lesson struct {
	client    *Client
	ID        int
	TeacherID int
	SubjectID int
	ClassID   int // 0 when Librus doesn't say
}

// Teacher returns the user who teaches the lesson.
func (l Lesson) Teacher(ctx context.Context) (*User, error) {
	return l.client.User(ctx, l.TeacherID)
}

// Subject returns the subject of the lesson.
func (l Lesson) Subject(ctx context.Context) (*Subject, error) {
	return l.client.Subject(ctx, l.SubjectID)
}

type lessonJSON struct {
	ID      flexInt `json:"Id"`
	Teacher ref     `json:"Teacher"`
	Subject ref     `json:"Subject"`
	Class   ref     `json:"Class"`
}

func (c *Client) toLesson(j lessonJSON) *Lesson {
	return &Lesson{
		client:    c,
		ID:        int(j.ID),
		TeacherID: int(j.Teacher.ID),
		SubjectID: int(j.Subject.ID),
		ClassID:   int(j.Class.ID),
	}
}

// Lesson returns the lesson with the given id.
func (c *Client) Lesson(ctx context.Context, id int) (*Lesson, error) {
	var resp struct {
		Lesson lessonJSON `json:"Lesson"`
	}
	if err := c.getJSON(ctx, "Lessons/"+strconv.Itoa(id), &resp); err != nil {
		return nil, err
	}
	return c.toLesson(resp.Lesson), nil
}

// Lessons returns all lessons of the student's class.
func (c *Client) Lessons(ctx context.Context) ([]*Lesson, error) {
	var resp struct {
		Lessons []lessonJSON `json:"Lessons"`
	}
	if err := c.getJSON(ctx, "Lessons", &resp); err != nil {
		return nil, err
	}
	lessons := make([]*Lesson, len(resp.Lessons))
	for i, j := range resp.Lessons {
		lessons[i] = c.toLesson(j)
	}
	return lessons, nil
}
