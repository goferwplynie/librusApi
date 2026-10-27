package librus

import (
	"context"
	"sort"
	"strconv"
	"strings"
)

// Average is a weighted average of grades.
type Average struct {
	Value float64 // 0 when Count is 0
	Count int     // number of grades counted into the average
}

// SubjectAverage holds the averages of a subject.
type SubjectAverage struct {
	SubjectID int
	Semester1 Average
	Semester2 Average
	// Year is the mean of the semester averages that have any grades, so it
	// equals the semester average until the other semester gets its first grade.
	Year Average
}

// Averages calculates the student's averages per subject from Grades and
// GradeCategories. Librus has its own Grades/Averages endpoint, but schools can
// disable it, so the average is computed here from grade weights.
//
// See CalculateAverages for the rules.
func (c *Client) Averages(ctx context.Context) ([]SubjectAverage, error) {
	grades, err := c.Grades(ctx)
	if err != nil {
		return nil, err
	}
	categories, err := c.GradeCategories(ctx)
	if err != nil {
		return nil, err
	}
	return CalculateAverages(grades, categories), nil
}

// CalculateAverages computes weighted averages per subject, sorted by subject id.
//
//   - A grade counts if its category has CountToTheAverage set; the category's
//     Weight is the weight of the grade. A grade with an unknown category counts
//     with weight 1.
//   - Semester and year grades and their proposals never count.
//   - Marks 1-6 count as is; "+" adds 0.5 and "-" subtracts 0.25 (so "4+" is 4.5
//     and "5-" is 4.75). Other grades, including a lone "+" or "-", "np" and "bz",
//     don't count.
//   - Point grades don't count; they have their own categories and rules.
func CalculateAverages(grades []*Grade, categories []GradeCategory) []SubjectAverage {
	byID := make(map[int]GradeCategory, len(categories))
	for _, category := range categories {
		byID[category.ID] = category
	}

	bySubject := map[int]*[2]gradeSum{} // by subject, then by semester

	for _, g := range grades {
		if g.IsSemesterGrade || g.IsSemesterProposition || g.IsFinalGrade || g.IsFinalProposition {
			continue
		}
		if g.Semester != 1 && g.Semester != 2 {
			continue
		}
		value, ok := gradeValue(g.Grade)
		if !ok {
			continue
		}
		weight := 1
		if category, found := byID[g.CategoryID]; found {
			if !category.CountToTheAverage {
				continue
			}
			weight = category.Weight
		}
		if weight <= 0 {
			continue
		}

		s := bySubject[g.SubjectID]
		if s == nil {
			s = &[2]gradeSum{}
			bySubject[g.SubjectID] = s
		}
		s[g.Semester-1].weighted += value * float64(weight)
		s[g.Semester-1].weights += float64(weight)
		s[g.Semester-1].count++
	}

	averages := make([]SubjectAverage, 0, len(bySubject))
	for subjectID, s := range bySubject {
		avg := SubjectAverage{
			SubjectID: subjectID,
			Semester1: s[0].average(),
			Semester2: s[1].average(),
		}
		avg.Year = Average{Count: avg.Semester1.Count + avg.Semester2.Count}
		switch {
		case avg.Semester1.Count > 0 && avg.Semester2.Count > 0:
			avg.Year.Value = (avg.Semester1.Value + avg.Semester2.Value) / 2
		case avg.Semester1.Count > 0:
			avg.Year.Value = avg.Semester1.Value
		default:
			avg.Year.Value = avg.Semester2.Value
		}
		averages = append(averages, avg)
	}
	sort.Slice(averages, func(i, j int) bool { return averages[i].SubjectID < averages[j].SubjectID })
	return averages
}

// gradeSum accumulates weighted grades of one subject in one semester.
type gradeSum struct {
	weighted, weights float64
	count             int
}

func (s gradeSum) average() Average {
	if s.count == 0 {
		return Average{}
	}
	return Average{Value: s.weighted / s.weights, Count: s.count}
}

// gradeValue converts a mark like "4", "4+" or "5-" to a number.
func gradeValue(grade string) (float64, bool) {
	mark, modifier := grade, ""
	if strings.HasSuffix(grade, "+") || strings.HasSuffix(grade, "-") {
		mark, modifier = grade[:len(grade)-1], grade[len(grade)-1:]
	}
	n, err := strconv.Atoi(mark)
	if err != nil || n < 1 || n > 6 {
		return 0, false
	}
	value := float64(n)
	switch modifier {
	case "+":
		value += 0.5
	case "-":
		value -= 0.25
	}
	return value, true
}
