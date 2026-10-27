package librus

import (
	"context"
	"strconv"
	"time"
)

// HomeworkAssignment is a homework set by a teacher.
type HomeworkAssignment struct {
	client    *Client
	ID        int
	TeacherID int
	LessonID  int
	Topic     string // e.g. "Notatka"
	Text      string
	Date      time.Time // when the homework was set
	DueDate   time.Time

	StudentIDs         []int // students the homework was set for
	ReadByIDs          []int // students who opened the homework
	MarkedAsDoneByIDs  []int // students who marked the homework as done
	MustSendAttachFile bool  // the student has to upload a file
	SendFilePossible   bool
}

// Teacher returns the user who set the homework.
func (h HomeworkAssignment) Teacher(ctx context.Context) (*User, error) {
	return h.client.User(ctx, h.TeacherID)
}

type homeworkAssignmentJSON struct {
	ID                      flexInt   `json:"Id"`
	Teacher                 ref       `json:"Teacher"`
	Lesson                  ref       `json:"Lesson"`
	Topic                   string    `json:"Topic"`
	Text                    string    `json:"Text"`
	Date                    string    `json:"Date"`
	DueDate                 string    `json:"DueDate"`
	Student                 []ref     `json:"Student"`
	StudentsWhoRead         []flexInt `json:"StudentsWhoRead"`
	StudentsWhoMarkedAsDone []flexInt `json:"StudentsWhoMarkedAsDone"`
	MustSendAttachFile      bool      `json:"MustSendAttachFile"`
	SendFilePossible        bool      `json:"SendFilePossible"`
}

func (c *Client) toHomeworkAssignment(j homeworkAssignmentJSON) (*HomeworkAssignment, error) {
	date, err := parseDate(j.Date)
	if err != nil {
		return nil, err
	}
	dueDate, err := parseDate(j.DueDate)
	if err != nil {
		return nil, err
	}
	h := &HomeworkAssignment{
		client:             c,
		ID:                 int(j.ID),
		TeacherID:          int(j.Teacher.ID),
		LessonID:           int(j.Lesson.ID),
		Topic:              j.Topic,
		Text:               j.Text,
		Date:               date,
		DueDate:            dueDate,
		MustSendAttachFile: j.MustSendAttachFile,
		SendFilePossible:   j.SendFilePossible,
	}
	for _, s := range j.Student {
		h.StudentIDs = append(h.StudentIDs, int(s.ID))
	}
	for _, id := range j.StudentsWhoRead {
		h.ReadByIDs = append(h.ReadByIDs, int(id))
	}
	for _, id := range j.StudentsWhoMarkedAsDone {
		h.MarkedAsDoneByIDs = append(h.MarkedAsDoneByIDs, int(id))
	}
	return h, nil
}

// HomeworkAssignment returns the homework with the given id.
func (c *Client) HomeworkAssignment(ctx context.Context, id int) (*HomeworkAssignment, error) {
	var resp struct {
		HomeWorkAssignment homeworkAssignmentJSON `json:"HomeWorkAssignment"`
	}
	if err := c.getJSON(ctx, "HomeWorkAssignments/"+strconv.Itoa(id), &resp); err != nil {
		return nil, err
	}
	return c.toHomeworkAssignment(resp.HomeWorkAssignment)
}

// HomeworkAssignments returns the homework set for the student.
func (c *Client) HomeworkAssignments(ctx context.Context) ([]*HomeworkAssignment, error) {
	var resp struct {
		HomeWorkAssignments []homeworkAssignmentJSON `json:"HomeWorkAssignments"`
	}
	if err := c.getJSON(ctx, "HomeWorkAssignments", &resp); err != nil {
		return nil, err
	}
	assignments := make([]*HomeworkAssignment, 0, len(resp.HomeWorkAssignments))
	for _, j := range resp.HomeWorkAssignments {
		h, err := c.toHomeworkAssignment(j)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, h)
	}
	return assignments, nil
}
