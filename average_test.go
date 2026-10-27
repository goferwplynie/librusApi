package librus

import (
	"math"
	"testing"
)

func TestGradeValue(t *testing.T) {
	tests := []struct {
		grade string
		want  float64
		ok    bool
	}{
		{"1", 1, true},
		{"6", 6, true},
		{"4+", 4.5, true},
		{"5-", 4.75, true},
		{"+", 0, false},
		{"-", 0, false},
		{"np", 0, false},
		{"bz", 0, false},
		{"", 0, false},
		{"7", 0, false},
		{"0", 0, false},
		{"6+", 6.5, true},
	}
	for _, tt := range tests {
		got, ok := gradeValue(tt.grade)
		if got != tt.want || ok != tt.ok {
			t.Errorf("gradeValue(%q) = %v, %v; want %v, %v", tt.grade, got, ok, tt.want, tt.ok)
		}
	}
}

func TestCalculateAverages(t *testing.T) {
	categories := []GradeCategory{
		{ID: 1, Weight: 1, CountToTheAverage: true},
		{ID: 2, Weight: 3, CountToTheAverage: true},
		{ID: 3, Weight: 5, CountToTheAverage: false},
	}
	grades := []*Grade{
		// subject 10, semester 1: (5*1 + 3*3 + 4+*1) / 5 = (5+9+4.5)/5 = 3.7
		{SubjectID: 10, Semester: 1, CategoryID: 1, Grade: "5"},
		{SubjectID: 10, Semester: 1, CategoryID: 2, Grade: "3"},
		{SubjectID: 10, Semester: 1, CategoryID: 1, Grade: "4+"},
		// ignored: category doesn't count, lone plus, semester grade, proposal, final grade, np
		{SubjectID: 10, Semester: 1, CategoryID: 3, Grade: "1"},
		{SubjectID: 10, Semester: 1, CategoryID: 1, Grade: "+"},
		{SubjectID: 10, Semester: 1, CategoryID: 1, Grade: "6", IsSemesterGrade: true},
		{SubjectID: 10, Semester: 1, CategoryID: 1, Grade: "6", IsSemesterProposition: true},
		{SubjectID: 10, Semester: 2, CategoryID: 1, Grade: "6", IsFinalGrade: true},
		{SubjectID: 10, Semester: 1, CategoryID: 1, Grade: "np"},
		// subject 20, both semesters: sem1 = 2, sem2 = 4, year = 3; unknown category counts with weight 1
		{SubjectID: 20, Semester: 1, CategoryID: 99, Grade: "2"},
		{SubjectID: 20, Semester: 2, CategoryID: 1, Grade: "4"},
		// subject 30 has no countable grades, so it isn't listed
		{SubjectID: 30, Semester: 1, CategoryID: 1, Grade: "+"},
	}

	got := CalculateAverages(grades, categories)
	if len(got) != 2 {
		t.Fatalf("got %d subjects, want 2: %+v", len(got), got)
	}

	near := func(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

	a := got[0]
	if a.SubjectID != 10 || !near(a.Semester1.Value, 3.7) || a.Semester1.Count != 3 || a.Semester2.Count != 0 {
		t.Errorf("subject 10 = %+v", a)
	}
	if !near(a.Year.Value, 3.7) || a.Year.Count != 3 {
		t.Errorf("subject 10 year = %+v, want it to equal semester 1", a.Year)
	}

	b := got[1]
	if b.SubjectID != 20 || !near(b.Semester1.Value, 2) || !near(b.Semester2.Value, 4) || !near(b.Year.Value, 3) || b.Year.Count != 2 {
		t.Errorf("subject 20 = %+v", b)
	}
}
