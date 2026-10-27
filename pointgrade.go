package librus

import (
	"context"
	"strconv"
	"time"
)

// PointGrade is a grade given as points (e.g. 8/10) instead of a mark.
type PointGrade struct {
	client     *Client
	ID         int
	LessonID   int
	SubjectID  int
	StudentID  int
	CategoryID int
	AddedByID  int
	CommentIDs []int

	Grade string  // as shown in Librus, e.g. "8.00"
	Value float64 // Grade as a number, 0 if it isn't numeric
	Date  time.Time
	// AddDate is when the teacher entered the grade.
	AddDate  time.Time
	Semester int
}

// Teacher returns the user who added the grade.
func (g PointGrade) Teacher(ctx context.Context) (*User, error) {
	return g.client.User(ctx, g.AddedByID)
}

// Subject returns the subject the grade was given for.
func (g PointGrade) Subject(ctx context.Context) (*Subject, error) {
	return g.client.Subject(ctx, g.SubjectID)
}

// Category returns the category of the grade, which holds the maximum number of points.
func (g PointGrade) Category(ctx context.Context) (*PointGradeCategory, error) {
	return g.client.PointGradeCategory(ctx, g.CategoryID)
}

// Comments returns the teacher's comments to the grade, one request per comment.
func (g PointGrade) Comments(ctx context.Context) ([]PointGradeComment, error) {
	var comments []PointGradeComment
	for _, id := range g.CommentIDs {
		comment, err := g.client.PointGradeComment(ctx, id)
		if err != nil {
			return nil, err
		}
		comments = append(comments, *comment)
	}
	return comments, nil
}

type pointGradeJSON struct {
	ID         flexInt `json:"Id"`
	Lesson     ref     `json:"Lesson"`
	Subject    ref     `json:"Subject"`
	Student    ref     `json:"Student"`
	Category   ref     `json:"Category"`
	AddedBy    ref     `json:"AddedBy"`
	Grade      string  `json:"Grade"`
	GradeValue string  `json:"GradeValue"`
	Date       string  `json:"Date"`
	AddDate    string  `json:"AddDate"`
	Semester   flexInt `json:"Semester"`
	Comments   []ref   `json:"Comments"`
}

func (c *Client) toPointGrade(j pointGradeJSON) (*PointGrade, error) {
	date, err := parseDate(j.Date)
	if err != nil {
		return nil, err
	}
	addDate, err := parseDateTime(j.AddDate)
	if err != nil {
		return nil, err
	}
	var commentIDs []int
	for _, comment := range j.Comments {
		commentIDs = append(commentIDs, int(comment.ID))
	}
	value, _ := strconv.ParseFloat(j.GradeValue, 64) // non-numeric grades stay 0

	return &PointGrade{
		client:     c,
		ID:         int(j.ID),
		LessonID:   int(j.Lesson.ID),
		SubjectID:  int(j.Subject.ID),
		StudentID:  int(j.Student.ID),
		CategoryID: int(j.Category.ID),
		AddedByID:  int(j.AddedBy.ID),
		CommentIDs: commentIDs,
		Grade:      j.Grade,
		Value:      value,
		Date:       date,
		AddDate:    addDate,
		Semester:   int(j.Semester),
	}, nil
}

// PointGrade returns the point grade with the given id.
func (c *Client) PointGrade(ctx context.Context, id int) (*PointGrade, error) {
	var resp struct {
		Grade pointGradeJSON `json:"Grade"`
	}
	if err := c.getJSON(ctx, "PointGrades/"+strconv.Itoa(id), &resp); err != nil {
		return nil, err
	}
	return c.toPointGrade(resp.Grade)
}

// PointGrades returns all point grades of the logged in student.
func (c *Client) PointGrades(ctx context.Context) ([]*PointGrade, error) {
	var resp struct {
		Grades []pointGradeJSON `json:"Grades"`
	}
	if err := c.getJSON(ctx, "PointGrades", &resp); err != nil {
		return nil, err
	}
	grades := make([]*PointGrade, 0, len(resp.Grades))
	for _, j := range resp.Grades {
		grade, err := c.toPointGrade(j)
		if err != nil {
			return nil, err
		}
		grades = append(grades, grade)
	}
	return grades, nil
}

// PointGradeCategory is a category of point grades. Unlike grade categories,
// every teacher has their own, each with its own point range.
type PointGradeCategory struct {
	ID                  int
	Name                string
	Short               string
	MinPoints           int
	MaxPoints           int
	Weight              int
	CountToTheAverage   bool
	ObligationToPerform bool
	TeacherID           int
	ColorID             int
}

type pointGradeCategoryJSON struct {
	ID                  flexInt `json:"Id"`
	Name                string  `json:"Name"`
	Short               string  `json:"Short"`
	ValueFrom           flexInt `json:"ValueFrom"`
	ValueTo             flexInt `json:"ValueTo"`
	Weight              flexInt `json:"Weight"`
	CountToTheAverage   bool    `json:"CountToTheAverage"`
	ObligationToPerform bool    `json:"ObligationToPerform"`
	Teacher             ref     `json:"Teacher"`
	Color               ref     `json:"Color"`
}

func (j pointGradeCategoryJSON) toCategory() PointGradeCategory {
	return PointGradeCategory{
		ID:                  int(j.ID),
		Name:                j.Name,
		Short:               j.Short,
		MinPoints:           int(j.ValueFrom),
		MaxPoints:           int(j.ValueTo),
		Weight:              int(j.Weight),
		CountToTheAverage:   j.CountToTheAverage,
		ObligationToPerform: j.ObligationToPerform,
		TeacherID:           int(j.Teacher.ID),
		ColorID:             int(j.Color.ID),
	}
}

// PointGradeCategory returns the point grade category with the given id.
func (c *Client) PointGradeCategory(ctx context.Context, id int) (*PointGradeCategory, error) {
	var resp struct {
		Category pointGradeCategoryJSON `json:"Category"`
	}
	if err := c.getJSON(ctx, "PointGrades/Categories/"+strconv.Itoa(id), &resp); err != nil {
		return nil, err
	}
	category := resp.Category.toCategory()
	return &category, nil
}

// PointGradeCategories returns all point grade categories.
func (c *Client) PointGradeCategories(ctx context.Context) ([]PointGradeCategory, error) {
	var resp struct {
		Categories []pointGradeCategoryJSON `json:"Categories"`
	}
	if err := c.getJSON(ctx, "PointGrades/Categories", &resp); err != nil {
		return nil, err
	}
	categories := make([]PointGradeCategory, len(resp.Categories))
	for i, j := range resp.Categories {
		categories[i] = j.toCategory()
	}
	return categories, nil
}

// PointGradeComment is a teacher's comment to a point grade.
type PointGradeComment struct {
	ID        int
	GradeID   int
	AddedByID int
	Text      string
}

// PointGradeComment returns the point grade comment with the given id.
func (c *Client) PointGradeComment(ctx context.Context, id int) (*PointGradeComment, error) {
	var resp struct {
		Comment gradeCommentJSON `json:"Comment"`
	}
	if err := c.getJSON(ctx, "PointGrades/Comments/"+strconv.Itoa(id), &resp); err != nil {
		return nil, err
	}
	return toPointGradeComment(resp.Comment), nil
}

// PointGradeComments returns all comments to the student's point grades.
func (c *Client) PointGradeComments(ctx context.Context) ([]PointGradeComment, error) {
	var resp struct {
		Comments []gradeCommentJSON `json:"Comments"`
	}
	if err := c.getJSON(ctx, "PointGrades/Comments", &resp); err != nil {
		return nil, err
	}
	comments := make([]PointGradeComment, len(resp.Comments))
	for i, j := range resp.Comments {
		comments[i] = *toPointGradeComment(j)
	}
	return comments, nil
}

func toPointGradeComment(j gradeCommentJSON) *PointGradeComment {
	comment := j.toComment()
	return &PointGradeComment{
		ID:        comment.ID,
		GradeID:   comment.GradeID,
		AddedByID: comment.AddedByID,
		Text:      comment.Text,
	}
}
