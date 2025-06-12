package timetable

type TimetableResponse struct {
	Timetable Timetable `json:"TimetableEntries"`
}

type Timetable struct {
	Day [][]Lesson `json:"day,omitempty"`
}

type Lesson struct {
	LessonNo            string  `json:"lesson_no,omitempty"`
	DayNo               string  `json:"day_no,omitempty"`
	Subject             Subject `json:"subject"`
	Teacher             Teacher `json:"teacher"`
	IsSubstitutionClass bool    `json:"is_substitution_class,omitempty"`
	IsCanceled          bool    `json:"is_canceled,omitempty"`
	SubstitutionNote    *string `json:"substitution_note,omitempty"`
	HourFrom            string  `json:"hour_from,omitempty"`
	HourTo              string  `json:"hour_to,omitempty"`
}

type Subject struct {
	Name  string `json:"name,omitempty"`
	Short string `json:"short,omitempty"`
}

type Teacher struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}
