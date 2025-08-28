package librusApi

type classResponse struct {
	Class Class  `json:"Class"`
	URL   string `json:"Url"`
}

type Class struct {
	EndFirstSemester string       `json:"EndFirstSemester"`
	ClassTutors      []ClassTutor `json:"ClassTutors"`
	Number           int64        `json:"Number"`
	ClassTutor       ClassTutor   `json:"ClassTutor"`
	Symbol           string       `json:"Symbol"`
	EndSchoolYear    string       `json:"EndSchoolYear"`
	ID               int64        `json:"Id"`
	Unit             ClassTutor   `json:"Unit"`
	BeginSchoolYear  string       `json:"BeginSchoolYear"`
}

type ClassTutor struct {
	ID  int64  `json:"Id"`
	URL string `json:"Url"`
}
