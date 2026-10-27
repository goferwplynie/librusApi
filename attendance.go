package librus

import (
	"context"
	"strconv"
	"time"
)

// Attendance is the student's attendance at a single lesson.
type Attendance struct {
	client *Client
	// ID is a number for regular lessons and a string like "t209235" for
	// lessons during trips.
	ID        string
	LessonID  int
	StudentID int
	TypeID    int // see Client.AttendanceTypes
	TripID    int // 0 if the lesson wasn't during a trip
	AddedByID int
	Date      time.Time
	AddDate   time.Time
	LessonNo  int
	Semester  int
}

// Type returns the kind of the attendance (presence, absence, lateness...).
func (a Attendance) Type(ctx context.Context) (*AttendanceType, error) {
	return a.client.AttendanceType(ctx, a.TypeID)
}

// Lesson returns the lesson the attendance was checked at.
func (a Attendance) Lesson(ctx context.Context) (*Lesson, error) {
	return a.client.Lesson(ctx, a.LessonID)
}

type attendanceJSON struct {
	ID       flexString `json:"Id"`
	Lesson   ref        `json:"Lesson"`
	Student  ref        `json:"Student"`
	Type     ref        `json:"Type"`
	Trip     ref        `json:"Trip"`
	AddedBy  ref        `json:"AddedBy"`
	Date     string     `json:"Date"`
	AddDate  string     `json:"AddDate"`
	LessonNo flexInt    `json:"LessonNo"`
	Semester flexInt    `json:"Semester"`
}

// Attendances returns the attendance of the student at all lessons, in the
// order given by Librus.
func (c *Client) Attendances(ctx context.Context) ([]*Attendance, error) {
	var resp struct {
		Attendances []attendanceJSON `json:"Attendances"`
	}
	if err := c.getJSON(ctx, "Attendances", &resp); err != nil {
		return nil, err
	}
	attendances := make([]*Attendance, 0, len(resp.Attendances))
	for _, j := range resp.Attendances {
		date, err := parseDate(j.Date)
		if err != nil {
			return nil, err
		}
		addDate, err := parseDateTime(j.AddDate)
		if err != nil {
			return nil, err
		}
		attendances = append(attendances, &Attendance{
			client:    c,
			ID:        string(j.ID),
			LessonID:  int(j.Lesson.ID),
			StudentID: int(j.Student.ID),
			TypeID:    int(j.Type.ID),
			TripID:    int(j.Trip.ID),
			AddedByID: int(j.AddedBy.ID),
			Date:      date,
			AddDate:   addDate,
			LessonNo:  int(j.LessonNo),
			Semester:  int(j.Semester),
		})
	}
	return attendances, nil
}

// AttendanceType is a kind of attendance, e.g. "Obecność" or "Nieobecność".
type AttendanceType struct {
	ID    int
	Name  string
	Short string // e.g. "nb"
	// IsPresence is false for absences (excused or not) and true for everything
	// that counts as being at the lesson (presence, lateness, being released).
	IsPresence bool
	Standard   bool
	ColorRGB   string // hex color without '#', e.g. "C6E0FF"
	Order      int
}

type attendanceTypeJSON struct {
	ID             flexInt `json:"Id"`
	Name           string  `json:"Name"`
	Short          string  `json:"Short"`
	IsPresenceKind bool    `json:"IsPresenceKind"`
	Standard       bool    `json:"Standard"`
	ColorRGB       string  `json:"ColorRGB"`
	Order          flexInt `json:"Order"`
}

func (j attendanceTypeJSON) toType() AttendanceType {
	return AttendanceType{
		ID:         int(j.ID),
		Name:       j.Name,
		Short:      j.Short,
		IsPresence: j.IsPresenceKind,
		Standard:   j.Standard,
		ColorRGB:   j.ColorRGB,
		Order:      int(j.Order),
	}
}

// AttendanceType returns the attendance type with the given id.
func (c *Client) AttendanceType(ctx context.Context, id int) (*AttendanceType, error) {
	var resp struct {
		Type attendanceTypeJSON `json:"Type"`
	}
	if err := c.getJSON(ctx, "Attendances/Types/"+strconv.Itoa(id), &resp); err != nil {
		return nil, err
	}
	t := resp.Type.toType()
	return &t, nil
}

// AttendanceTypes returns all kinds of attendance, sorted by Librus's order.
func (c *Client) AttendanceTypes(ctx context.Context) ([]AttendanceType, error) {
	var resp struct {
		Types []attendanceTypeJSON `json:"Types"`
	}
	if err := c.getJSON(ctx, "Attendances/Types", &resp); err != nil {
		return nil, err
	}
	types := make([]AttendanceType, len(resp.Types))
	for i, j := range resp.Types {
		types[i] = j.toType()
	}
	return types, nil
}
