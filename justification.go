package librus

import (
	"context"
	"strconv"
	"time"
)

// Justification is an absence justification sent by a parent.
type Justification struct {
	ID      int
	Message string // message from the parent to the teachers
	// Status is Librus's status of the justification, e.g. "accept".
	Status   string
	PostDate time.Time
	DateFrom time.Time
	DateTo   time.Time
	// JustifiedAbsences is the number of absences that were justified.
	JustifiedAbsences int
	HasAttachment     bool
	// NotifiedTeachers are the names of teachers who got the justification, as
	// "LastName FirstName".
	NotifiedTeachers []string
}

type justificationJSON struct {
	ID                  flexInt `json:"id"`
	MessageFromParent   string  `json:"messageFromParent"`
	PostDate            string  `json:"postDate"`
	JustificationStatus string  `json:"justificationStatus"`
	Attachment          bool    `json:"attachment"`
	DateFrom            string  `json:"dateFrom"`
	DateTo              string  `json:"dateTo"`
	JustifiedAbsences   flexInt `json:"justifiedAbsences"`
	NotifiedTeachers    []struct {
		Name string `json:"name"`
	} `json:"notifiedTeachers"`
}

func (j justificationJSON) toJustification() (*Justification, error) {
	postDate, err := parseDateTime(j.PostDate)
	if err != nil {
		return nil, err
	}
	from, err := parseDate(j.DateFrom)
	if err != nil {
		return nil, err
	}
	to, err := parseDate(j.DateTo)
	if err != nil {
		return nil, err
	}
	justification := &Justification{
		ID:                int(j.ID),
		Message:           j.MessageFromParent,
		Status:            j.JustificationStatus,
		PostDate:          postDate,
		DateFrom:          from,
		DateTo:            to,
		JustifiedAbsences: int(j.JustifiedAbsences),
		HasAttachment:     j.Attachment,
	}
	for _, t := range j.NotifiedTeachers {
		justification.NotifiedTeachers = append(justification.NotifiedTeachers, t.Name)
	}
	return justification, nil
}

// Justification returns the justification with the given id.
func (c *Client) Justification(ctx context.Context, id int) (*Justification, error) {
	var resp struct {
		Data justificationJSON `json:"data"`
	}
	if err := c.getJSON(ctx, "Justifications/"+strconv.Itoa(id), &resp); err != nil {
		return nil, err
	}
	return resp.Data.toJustification()
}

// Justifications returns the absence justifications sent for the student.
func (c *Client) Justifications(ctx context.Context) ([]*Justification, error) {
	var resp struct {
		Data []justificationJSON `json:"data"`
	}
	if err := c.getJSON(ctx, "Justifications", &resp); err != nil {
		return nil, err
	}
	justifications := make([]*Justification, 0, len(resp.Data))
	for _, j := range resp.Data {
		justification, err := j.toJustification()
		if err != nil {
			return nil, err
		}
		justifications = append(justifications, justification)
	}
	return justifications, nil
}
