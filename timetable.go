package librus

import (
	"context"
	"sort"
	"time"
)

// TimetableLesson is a single lesson in the timetable.
type TimetableLesson struct {
	Date     time.Time // day of the lesson, midnight in Librus's time zone
	LessonNo int       // lesson number in the day, as shown in Librus
	Start    time.Time
	End      time.Time

	Subject     TimetableSubject
	Teacher     TimetableTeacher
	ClassroomID int    // 0 if no classroom is assigned; see Client.Classrooms
	GroupName   string // name of the group (e.g. a language group), empty for lessons of the whole class

	IsCanceled bool
	// IsSubstitution is set when the lesson replaces another one; Original then describes the replaced lesson.
	IsSubstitution bool
	// Note is the teacher's note for a canceled or substituted lesson.
	Note     string
	Original *OriginalLesson
}

type TimetableSubject struct {
	ID    int
	Name  string
	Short string
}

type TimetableTeacher struct {
	ID        int
	FirstName string
	LastName  string
}

// OriginalLesson is the lesson that a substitution replaced.
type OriginalLesson struct {
	client      *Client
	Date        time.Time
	LessonNo    int
	Start       time.Time
	End         time.Time
	SubjectID   int
	TeacherID   int
	ClassroomID int
}

// Subject fetches the subject of the replaced lesson. It costs one request,
// so it's only made when called.
func (o OriginalLesson) Subject(ctx context.Context) (*Subject, error) {
	return o.client.Subject(ctx, o.SubjectID)
}

type timetableResponse struct {
	// date -> lesson slots -> lessons in the slot (usually 0 or 1, more for groups)
	Timetable map[string][][]timetableLessonJSON `json:"Timetable"`
}

type timetableLessonJSON struct {
	LessonNo flexInt `json:"LessonNo"`
	HourFrom string  `json:"HourFrom"`
	HourTo   string  `json:"HourTo"`
	Subject  struct {
		ID    flexInt `json:"Id"`
		Name  string  `json:"Name"`
		Short string  `json:"Short"`
	} `json:"Subject"`
	Teacher struct {
		ID        flexInt `json:"Id"`
		FirstName string  `json:"FirstName"`
		LastName  string  `json:"LastName"`
	} `json:"Teacher"`
	Classroom           ref     `json:"Classroom"`
	VirtualClassName    string  `json:"VirtualClassName"`
	IsCanceled          bool    `json:"IsCanceled"`
	IsSubstitutionClass bool    `json:"IsSubstitutionClass"`
	SubstitutionNote    *string `json:"SubstitutionNote"`

	OrgDate      string  `json:"OrgDate"`
	OrgLessonNo  flexInt `json:"OrgLessonNo"`
	OrgHourFrom  string  `json:"OrgHourFrom"`
	OrgHourTo    string  `json:"OrgHourTo"`
	OrgSubject   ref     `json:"OrgSubject"`
	OrgTeacher   ref     `json:"OrgTeacher"`
	OrgClassroom ref     `json:"OrgClassroom"`
}

// Timetable returns the lessons of the week (Monday to Sunday) containing day,
// sorted by date and lesson number.
func (c *Client) Timetable(ctx context.Context, day time.Time) ([]TimetableLesson, error) {
	day = day.In(librusLocation)
	monday := day.AddDate(0, 0, -(int(day.Weekday())+6)%7)
	return c.timetable(ctx, "Timetables?weekStart="+monday.Format(time.DateOnly))
}

// TimetableDay returns the lessons of a single day, sorted by lesson number.
func (c *Client) TimetableDay(ctx context.Context, day time.Time) ([]TimetableLesson, error) {
	return c.timetable(ctx, "Timetables?day="+day.In(librusLocation).Format(time.DateOnly))
}

func (c *Client) timetable(ctx context.Context, resource string) ([]TimetableLesson, error) {
	var resp timetableResponse
	if err := c.getJSON(ctx, resource, &resp); err != nil {
		return nil, err
	}

	var lessons []TimetableLesson
	for dateStr, slots := range resp.Timetable {
		date, err := parseDate(dateStr)
		if err != nil {
			return nil, err
		}
		for _, slot := range slots {
			for _, l := range slot {
				lesson, err := l.toLesson(c, date)
				if err != nil {
					return nil, err
				}
				lessons = append(lessons, lesson)
			}
		}
	}

	sort.Slice(lessons, func(i, j int) bool {
		if !lessons[i].Date.Equal(lessons[j].Date) {
			return lessons[i].Date.Before(lessons[j].Date)
		}
		return lessons[i].LessonNo < lessons[j].LessonNo
	})
	return lessons, nil
}

func (l timetableLessonJSON) toLesson(c *Client, date time.Time) (TimetableLesson, error) {
	start, err := atHour(date, l.HourFrom)
	if err != nil {
		return TimetableLesson{}, err
	}
	end, err := atHour(date, l.HourTo)
	if err != nil {
		return TimetableLesson{}, err
	}

	lesson := TimetableLesson{
		Date:     date,
		LessonNo: int(l.LessonNo),
		Start:    start,
		End:      end,
		Subject: TimetableSubject{
			ID:    int(l.Subject.ID),
			Name:  l.Subject.Name,
			Short: l.Subject.Short,
		},
		Teacher: TimetableTeacher{
			ID:        int(l.Teacher.ID),
			FirstName: l.Teacher.FirstName,
			LastName:  l.Teacher.LastName,
		},
		ClassroomID:    int(l.Classroom.ID),
		GroupName:      l.VirtualClassName,
		IsCanceled:     l.IsCanceled,
		IsSubstitution: l.IsSubstitutionClass,
	}
	if l.SubstitutionNote != nil {
		lesson.Note = *l.SubstitutionNote
	}

	if l.OrgDate != "" {
		orgDate, err := parseDate(l.OrgDate)
		if err != nil {
			return TimetableLesson{}, err
		}
		orgStart, err := atHour(orgDate, l.OrgHourFrom)
		if err != nil {
			return TimetableLesson{}, err
		}
		orgEnd, err := atHour(orgDate, l.OrgHourTo)
		if err != nil {
			return TimetableLesson{}, err
		}
		lesson.Original = &OriginalLesson{
			client:      c,
			Date:        orgDate,
			LessonNo:    int(l.OrgLessonNo),
			Start:       orgStart,
			End:         orgEnd,
			SubjectID:   int(l.OrgSubject.ID),
			TeacherID:   int(l.OrgTeacher.ID),
			ClassroomID: int(l.OrgClassroom.ID),
		}
	}
	return lesson, nil
}
