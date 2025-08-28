package librusApi

type classroomsResponse struct {
	Classrooms []Classroom `json:"Classrooms"`
	URL        string      `json:"Url"`
}

type Classroom struct {
	SchoolCommonRoom bool    `json:"SchoolCommonRoom"`
	Description      *string `json:"Description,omitempty"`
	Symbol           string  `json:"Symbol"`
	Size             int64   `json:"Size"`
	ID               int64   `json:"Id"`
	Name             string  `json:"Name"`
}
