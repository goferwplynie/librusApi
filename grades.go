package librusApi

import (
	"encoding/json"
	"fmt"
	"time"
)

func GetGrades(c *Client) ([]*Grade, error) {
	grades := make([]*Grade, 0)
	gradesJsonResp, err := c.Get("Grades", nil)
	if err != nil {
		return nil, err
	}
	defer gradesJsonResp.Body.Close()

	var gradesResp GradesResponse
	if err := json.NewDecoder(gradesJsonResp.Body).Decode(&gradesResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json:\n%w", err)
	}

	for _, v := range gradesResp.Grades {
		date, err := time.Parse(time.DateOnly, v.Date)
		if err != nil {
			return nil, fmt.Errorf("failed parsing time:\n%w", err)
		}
		addDate, err := time.Parse(time.DateTime, v.AddDate)
		if err != nil {
			return nil, fmt.Errorf("failed parsing time:\n%w", err)
		}

		grade := newGrade(GradeParams{
			client:             c,
			SubjectId:          v.Subject.ID,
			StudentId:          v.Student.ID,
			CategoryId:         v.Category.ID,
			Semester:           v.Semester,
			AddedById:          v.AddedBy.ID,
			Grade:              v.Grade,
			Date:               date,
			AddDate:            addDate,
			IsConstituent:      v.IsConstituent,
			IsSemesterGrade:    v.IsSemester,
			IsFinalGrade:       v.IsFinal,
			IsFinalProposition: v.IsFinalProposition,
		})

		grades = append(grades, grade)

	}

	return grades, nil
}

func newGrade(GradeParams GradeParams) *Grade {
	return &Grade{
		client:                GradeParams.client,
		SubjectId:             GradeParams.SubjectId,
		StudentId:             GradeParams.StudentId,
		CategoryId:            GradeParams.CategoryId,
		Semester:              GradeParams.Semester,
		AddedById:             GradeParams.AddedById,
		Grade:                 GradeParams.Grade,
		Date:                  GradeParams.Date,
		AddDate:               GradeParams.AddDate,
		IsConstituent:         GradeParams.IsConstituent,
		IsSemesterGrade:       GradeParams.IsSemesterGrade,
		IsSemesterProposition: GradeParams.IsSemesterProposition,
		IsFinalGrade:          GradeParams.IsFinalGrade,
		IsFinalProposition:    GradeParams.IsFinalProposition,
	}
}
