package librus

import (
	"context"
	"sort"
	"strconv"
	"strings"
)

// Average is an average of grades.
type Average struct {
	// Value is the weighted average of marks, or the percent of possible points
	// (0-100) for a subject graded with points, see SubjectAverage.UsesPoints.
	// It is 0 when Count is 0.
	Value float64
	Count int // number of grades counted into the average
}

// SubjectAverage holds the averages of a subject.
type SubjectAverage struct {
	SubjectID int
	// UsesPoints is true for subjects graded with points instead of marks, i.e.
	// those that have point grades and no marks that count. Their averages are
	// the percent of possible points scored: all scored points divided by all
	// maximum points, so 8/10 and 3/20 give 11/30 = 36.7.
	UsesPoints bool
	Semester1  Average
	Semester2  Average
	// Year is the mean of the semester averages that have any grades (so it
	// equals the semester average until the other semester gets its first
	// grade). For subjects graded with points the points of both semesters are
	// added together instead.
	Year Average
}

// Averages calculates the student's averages per subject from Grades,
// GradeCategories, PointGrades and PointGradeCategories. Librus has its own
// Grades/Averages endpoint, but schools can disable it, so the average is
// computed here from grade weights.
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
	pointGrades, err := c.PointGrades(ctx)
	if err != nil {
		return nil, err
	}
	pointCategories, err := c.PointGradeCategories(ctx)
	if err != nil {
		return nil, err
	}
	return CalculateAverages(grades, categories, pointGrades, pointCategories), nil
}

// CalculateAverages computes averages per subject, sorted by subject id.
//
// Marks:
//   - A grade counts if its category has CountToTheAverage set; the category's
//     Weight is the weight of the grade. A grade with an unknown category counts
//     with weight 1.
//   - Semester and year grades and their proposals never count.
//   - Only the mark counts: "4+" and "4-" are both 4. Other grades, including a
//     lone "+" or "-", "np" and "bz", don't count.
//
// Point grades:
//   - A subject with point grades and no counted marks is graded with points
//     (SubjectAverage.UsesPoints). Its average is all scored points divided by
//     all maximum points (MaxPoints of the grade's category), in percent.
//     Weights are not applied.
//   - Grades that aren't numeric and grades whose category is unknown or has no
//     point range are skipped.
//   - A subject that has both counted marks and point grades is averaged from
//     the marks only.
func CalculateAverages(grades []*Grade, categories []GradeCategory, pointGrades []*PointGrade, pointCategories []PointGradeCategory) []SubjectAverage {
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
		s[g.Semester-1].add(value*float64(weight), float64(weight))
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

	for subjectID, s := range pointSums(pointGrades, pointCategories) {
		if _, hasMarks := bySubject[subjectID]; hasMarks {
			continue
		}
		year := gradeSum{
			weighted: s[0].weighted + s[1].weighted,
			weights:  s[0].weights + s[1].weights,
			count:    s[0].count + s[1].count,
		}
		averages = append(averages, SubjectAverage{
			SubjectID:  subjectID,
			UsesPoints: true,
			Semester1:  s[0].percent(),
			Semester2:  s[1].percent(),
			Year:       year.percent(),
		})
	}

	sort.Slice(averages, func(i, j int) bool { return averages[i].SubjectID < averages[j].SubjectID })
	return averages
}

// pointSums adds up scored points (weighted) and maximum points (weights) of
// each subject by semester.
func pointSums(grades []*PointGrade, categories []PointGradeCategory) map[int]*[2]gradeSum {
	byID := make(map[int]PointGradeCategory, len(categories))
	for _, category := range categories {
		byID[category.ID] = category
	}

	bySubject := map[int]*[2]gradeSum{}
	for _, g := range grades {
		if g.Semester != 1 && g.Semester != 2 {
			continue
		}
		category, found := byID[g.CategoryID]
		if !found || category.MaxPoints <= category.MinPoints {
			continue
		}
		points, err := strconv.ParseFloat(g.Grade, 64)
		if err != nil {
			continue
		}

		s := bySubject[g.SubjectID]
		if s == nil {
			s = &[2]gradeSum{}
			bySubject[g.SubjectID] = s
		}
		s[g.Semester-1].add(points, float64(category.MaxPoints))
	}
	return bySubject
}

// gradeSum accumulates grades of one subject in one semester: weighted is the
// sum of values times weights (or scored points), weights the sum of weights
// (or maximum points).
type gradeSum struct {
	weighted, weights float64
	count             int
}

func (s *gradeSum) add(weighted, weights float64) {
	s.weighted += weighted
	s.weights += weights
	s.count++
}

func (s gradeSum) average() Average {
	if s.count == 0 {
		return Average{}
	}
	return Average{Value: s.weighted / s.weights, Count: s.count}
}

// percent is average as a percentage, for sums of points.
func (s gradeSum) percent() Average {
	avg := s.average()
	avg.Value *= 100
	return avg
}

// gradeValue converts a mark like "4", "4+" or "5-" to a number. A "+" or "-"
// after the mark is ignored.
func gradeValue(grade string) (float64, bool) {
	mark := strings.TrimRight(grade, "+-")
	if len(grade)-len(mark) > 1 {
		return 0, false
	}
	n, err := strconv.Atoi(mark)
	if err != nil || n < 1 || n > 6 {
		return 0, false
	}
	return float64(n), true
}
