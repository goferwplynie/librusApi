package librus

import (
	"context"
	"strconv"
	"time"
)

// Event is an entry of the school calendar ("terminarz"): a test, a quiz, a
// trip, a long-term project, a parent-teacher meeting and so on. The kind is
// given by the category, see Client.EventCategories.
type Event struct {
	client      *Client
	ID          int
	Content     string
	Date        time.Time
	LessonNo    int       // 0 if the event isn't bound to a lesson
	Start       time.Time // zero if Librus gave no hour
	End         time.Time // zero if Librus gave no hour
	CategoryID  int
	CreatedByID int
	ClassID     int
	SubjectID   int // 0 if the event isn't bound to a subject
	AddDate     time.Time
}

// Category returns the category of the event.
func (e Event) Category(ctx context.Context) (*EventCategory, error) {
	return e.client.EventCategory(ctx, e.CategoryID)
}

// CreatedBy returns the user who added the event.
func (e Event) CreatedBy(ctx context.Context) (*User, error) {
	return e.client.User(ctx, e.CreatedByID)
}

// Subject returns the subject of the event. It returns nil without any request
// if the event isn't bound to a subject.
func (e Event) Subject(ctx context.Context) (*Subject, error) {
	if e.SubjectID == 0 {
		return nil, nil
	}
	return e.client.Subject(ctx, e.SubjectID)
}

type eventJSON struct {
	ID        flexInt `json:"Id"`
	Content   string  `json:"Content"`
	Date      string  `json:"Date"`
	Category  ref     `json:"Category"`
	LessonNo  flexInt `json:"LessonNo"`
	TimeFrom  string  `json:"TimeFrom"`
	TimeTo    string  `json:"TimeTo"`
	CreatedBy ref     `json:"CreatedBy"`
	Class     ref     `json:"Class"`
	Subject   ref     `json:"Subject"`
	AddDate   string  `json:"AddDate"`
}

func (c *Client) toEvent(j eventJSON) (*Event, error) {
	date, err := parseDate(j.Date)
	if err != nil {
		return nil, err
	}
	event := &Event{
		client:      c,
		ID:          int(j.ID),
		Content:     j.Content,
		Date:        date,
		LessonNo:    int(j.LessonNo),
		CategoryID:  int(j.Category.ID),
		CreatedByID: int(j.CreatedBy.ID),
		ClassID:     int(j.Class.ID),
		SubjectID:   int(j.Subject.ID),
	}
	if j.TimeFrom != "" {
		if event.Start, err = atHour(date, j.TimeFrom); err != nil {
			return nil, err
		}
	}
	if j.TimeTo != "" {
		if event.End, err = atHour(date, j.TimeTo); err != nil {
			return nil, err
		}
	}
	if j.AddDate != "" {
		if event.AddDate, err = parseDateTime(j.AddDate); err != nil {
			return nil, err
		}
	}
	return event, nil
}

// Event returns the calendar event with the given id.
func (c *Client) Event(ctx context.Context, id int) (*Event, error) {
	var resp struct {
		HomeWork eventJSON `json:"HomeWork"`
	}
	if err := c.getJSON(ctx, "HomeWorks/"+strconv.Itoa(id), &resp); err != nil {
		return nil, err
	}
	return c.toEvent(resp.HomeWork)
}

// Events returns all events of the student's class calendar (tests, quizzes,
// trips and so on), in the order given by Librus.
func (c *Client) Events(ctx context.Context) ([]*Event, error) {
	var resp struct {
		HomeWorks []eventJSON `json:"HomeWorks"`
	}
	if err := c.getJSON(ctx, "HomeWorks", &resp); err != nil {
		return nil, err
	}
	events := make([]*Event, 0, len(resp.HomeWorks))
	for _, j := range resp.HomeWorks {
		event, err := c.toEvent(j)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

// EventCategory is a kind of calendar event, e.g. "Sprawdzian".
type EventCategory struct {
	ID      int
	Name    string
	ColorID int
}

type eventCategoryJSON struct {
	ID    flexInt `json:"Id"`
	Name  string  `json:"Name"`
	Color ref     `json:"Color"`
}

// EventCategory returns the event category with the given id.
func (c *Client) EventCategory(ctx context.Context, id int) (*EventCategory, error) {
	var resp struct {
		Category eventCategoryJSON `json:"Category"`
	}
	if err := c.getJSON(ctx, "HomeWorks/Categories/"+strconv.Itoa(id), &resp); err != nil {
		return nil, err
	}
	return &EventCategory{ID: int(resp.Category.ID), Name: resp.Category.Name, ColorID: int(resp.Category.Color.ID)}, nil
}

// EventCategories returns all kinds of calendar events.
func (c *Client) EventCategories(ctx context.Context) ([]EventCategory, error) {
	var resp struct {
		Categories []eventCategoryJSON `json:"Categories"`
	}
	if err := c.getJSON(ctx, "HomeWorks/Categories", &resp); err != nil {
		return nil, err
	}
	categories := make([]EventCategory, len(resp.Categories))
	for i, j := range resp.Categories {
		categories[i] = EventCategory{ID: int(j.ID), Name: j.Name, ColorID: int(j.Color.ID)}
	}
	return categories, nil
}
