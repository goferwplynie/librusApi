package librus

import (
	"fmt"
	"time"
)

// librusLocation is the time zone of dates and hours returned by Librus. Falls
// back to the local time zone when the system has no tz database.
var librusLocation = func() *time.Location {
	loc, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		return time.Local
	}
	return loc
}()

// parseDate parses a "2006-01-02" date as midnight in Librus's time zone.
func parseDate(s string) (time.Time, error) {
	t, err := time.ParseInLocation(time.DateOnly, s, librusLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("librus: parsing date: %w", err)
	}
	return t, nil
}

// parseDateTime parses a "2006-01-02 15:04:05" timestamp in Librus's time zone.
func parseDateTime(s string) (time.Time, error) {
	t, err := time.ParseInLocation(time.DateTime, s, librusLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("librus: parsing date: %w", err)
	}
	return t, nil
}

// atHour returns date at the "15:04" or "15:04:05" hour.
func atHour(date time.Time, hour string) (time.Time, error) {
	layout := "15:04"
	if len(hour) > len(layout) {
		layout = time.TimeOnly
	}
	h, err := time.Parse(layout, hour)
	if err != nil {
		return time.Time{}, fmt.Errorf("librus: parsing hour: %w", err)
	}
	return time.Date(date.Year(), date.Month(), date.Day(), h.Hour(), h.Minute(), 0, 0, date.Location()), nil
}
