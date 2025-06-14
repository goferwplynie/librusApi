package grades

type GradesResponse struct {
	Grades []struct {
		ID     int `json:"Id"`
		Lesson struct {
			ID  int    `json:"Id"`
			URL string `json:"Url"`
		} `json:"Lesson"`
		Subject struct {
			ID  int    `json:"Id"`
			URL string `json:"Url"`
		} `json:"Subject"`
		Student struct {
			ID  int    `json:"Id"`
			URL string `json:"Url"`
		} `json:"Student"`
		Category struct {
			ID  int    `json:"Id"`
			URL string `json:"Url"`
		} `json:"Category"`
		AddedBy struct {
			ID  int    `json:"Id"`
			URL string `json:"Url"`
		} `json:"AddedBy"`
		Grade                 string `json:"Grade"`
		Date                  string `json:"Date"`
		AddDate               string `json:"AddDate"`
		Semester              int    `json:"Semester"`
		IsConstituent         bool   `json:"IsConstituent"`
		IsSemester            bool   `json:"IsSemester"`
		IsSemesterProposition bool   `json:"IsSemesterProposition"`
		IsFinal               bool   `json:"IsFinal"`
		IsFinalProposition    bool   `json:"IsFinalProposition"`
		Comments              []struct {
			ID  int    `json:"Id"`
			URL string `json:"Url"`
		} `json:"Comments,omitempty"`
	} `json:"Grades"`
}
