package librus

import (
	"context"
	"sort"
	"time"
)

// SchoolFreeDay is a day off (or a range of days) for the whole school: a
// holiday, a break or a principal's day off.
type SchoolFreeDay struct {
	ID       int
	Name     string
	DateFrom time.Time
	DateTo   time.Time // equal to DateFrom for a single day
	UnitIDs  []int     // school units the day off applies to
}

// Includes reports whether day falls in the range, ignoring the time of day.
func (d SchoolFreeDay) Includes(day time.Time) bool {
	day = day.In(librusLocation)
	day = time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, librusLocation)
	return !day.Before(d.DateFrom) && !day.After(d.DateTo)
}

// SchoolFreeDays returns the school's days off, sorted by date.
func (c *Client) SchoolFreeDays(ctx context.Context) ([]SchoolFreeDay, error) {
	var resp struct {
		SchoolFreeDays []struct {
			ID       flexInt `json:"Id"`
			Name     string  `json:"Name"`
			DateFrom string  `json:"DateFrom"`
			DateTo   string  `json:"DateTo"`
			Units    []ref   `json:"Units"`
		} `json:"SchoolFreeDays"`
	}
	if err := c.getJSON(ctx, "SchoolFreeDays", &resp); err != nil {
		return nil, err
	}

	days := make([]SchoolFreeDay, 0, len(resp.SchoolFreeDays))
	for _, j := range resp.SchoolFreeDays {
		from, err := parseDate(j.DateFrom)
		if err != nil {
			return nil, err
		}
		to, err := parseDate(j.DateTo)
		if err != nil {
			return nil, err
		}
		day := SchoolFreeDay{ID: int(j.ID), Name: j.Name, DateFrom: from, DateTo: to}
		for _, u := range j.Units {
			day.UnitIDs = append(day.UnitIDs, int(u.ID))
		}
		days = append(days, day)
	}
	sort.Slice(days, func(i, j int) bool { return days[i].DateFrom.Before(days[j].DateFrom) })
	return days, nil
}
