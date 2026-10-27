package librus

import (
	"context"
	"strconv"
)

// GradeCategory is a category of grades (e.g. "sprawdzian"), set by a teacher
// or by the school. Its weight and CountToTheAverage decide how grades of this
// category are counted into the average.
type GradeCategory struct {
	ID    int
	Name  string
	Short string // often empty
	// Weight is 0 for categories that don't count (semester and year grades).
	Weight            int
	CountToTheAverage bool
	Standard          bool
	BlockAnyGrades    bool

	// Kind of the category for semester and year grades.
	IsFinal               bool // year grade
	IsFinalProposition    bool // proposed year grade
	IsSemester            bool // semester grade
	IsSemesterProposition bool // proposed semester grade
	ParentCategoryID      int  // 0 if the category has no parent
	TeacherID             int  // 0 for categories defined by the school
	ColorID               int
	ObligationToPerform   bool
}

type gradeCategoryJSON struct {
	ID                     flexInt `json:"Id"`
	Name                   string  `json:"Name"`
	Short                  string  `json:"Short"`
	Weight                 flexInt `json:"Weight"`
	CountToTheAverage      bool    `json:"CountToTheAverage"`
	Standard               bool    `json:"Standard"`
	BlockAnyGrades         bool    `json:"BlockAnyGrades"`
	IsFinal                bool    `json:"IsFinal"`
	IsFinalProposition     bool    `json:"IsFinalProposition"`
	IsSemestral            bool    `json:"IsSemestral"`
	IsSemestralProposition bool    `json:"IsSemestralProposition"`
	ParentCategory         ref     `json:"ParentCategory"`
	Teacher                ref     `json:"Teacher"`
	Color                  ref     `json:"Color"`
	ObligationToPerform    bool    `json:"ObligationToPerform"`
}

func (j gradeCategoryJSON) toCategory() GradeCategory {
	return GradeCategory{
		ID:                    int(j.ID),
		Name:                  j.Name,
		Short:                 j.Short,
		Weight:                int(j.Weight),
		CountToTheAverage:     j.CountToTheAverage,
		Standard:              j.Standard,
		BlockAnyGrades:        j.BlockAnyGrades,
		IsFinal:               j.IsFinal,
		IsFinalProposition:    j.IsFinalProposition,
		IsSemester:            j.IsSemestral,
		IsSemesterProposition: j.IsSemestralProposition,
		ParentCategoryID:      int(j.ParentCategory.ID),
		TeacherID:             int(j.Teacher.ID),
		ColorID:               int(j.Color.ID),
		ObligationToPerform:   j.ObligationToPerform,
	}
}

// GradeCategory returns the grade category with the given id.
func (c *Client) GradeCategory(ctx context.Context, id int) (*GradeCategory, error) {
	var resp struct {
		Category gradeCategoryJSON `json:"Category"`
	}
	if err := c.getJSON(ctx, "Grades/Categories/"+strconv.Itoa(id), &resp); err != nil {
		return nil, err
	}
	category := resp.Category.toCategory()
	return &category, nil
}

// GradeCategories returns all grade categories visible to the student.
func (c *Client) GradeCategories(ctx context.Context) ([]GradeCategory, error) {
	var resp struct {
		Categories []gradeCategoryJSON `json:"Categories"`
	}
	if err := c.getJSON(ctx, "Grades/Categories", &resp); err != nil {
		return nil, err
	}
	categories := make([]GradeCategory, len(resp.Categories))
	for i, j := range resp.Categories {
		categories[i] = j.toCategory()
	}
	return categories, nil
}

// GradeComment is a teacher's comment to a grade.
type GradeComment struct {
	ID        int
	GradeID   int
	AddedByID int
	Text      string
}

type gradeCommentJSON struct {
	ID      flexInt `json:"Id"`
	AddedBy ref     `json:"AddedBy"`
	Grade   ref     `json:"Grade"`
	Text    string  `json:"Text"`
}

func (j gradeCommentJSON) toComment() GradeComment {
	return GradeComment{
		ID:        int(j.ID),
		GradeID:   int(j.Grade.ID),
		AddedByID: int(j.AddedBy.ID),
		Text:      j.Text,
	}
}

// GradeComment returns the grade comment with the given id.
func (c *Client) GradeComment(ctx context.Context, id int) (*GradeComment, error) {
	var resp struct {
		Comment gradeCommentJSON `json:"Comment"`
	}
	if err := c.getJSON(ctx, "Grades/Comments/"+strconv.Itoa(id), &resp); err != nil {
		return nil, err
	}
	comment := resp.Comment.toComment()
	return &comment, nil
}

// GradeComments returns all comments to the student's grades.
func (c *Client) GradeComments(ctx context.Context) ([]GradeComment, error) {
	var resp struct {
		Comments []gradeCommentJSON `json:"Comments"`
	}
	if err := c.getJSON(ctx, "Grades/Comments", &resp); err != nil {
		return nil, err
	}
	comments := make([]GradeComment, len(resp.Comments))
	for i, j := range resp.Comments {
		comments[i] = j.toComment()
	}
	return comments, nil
}
