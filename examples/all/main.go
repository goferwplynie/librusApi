// Shows everything the library can fetch for a student account. Each section
// prints its own error and the program carries on, because schools can switch
// off some of the features.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	librus "github.com/goferwplynie/librusApi"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	ctx := context.Background()
	client := librus.New()
	if err := client.Login(ctx, os.Getenv("LIBRUS_LOGIN"), os.Getenv("LIBRUS_PASSWORD")); err != nil {
		panic(err)
	}

	// Lists are fetched once and used to turn ids into names.
	subjects := map[int]string{}
	if list, err := client.Subjects(ctx); check("subjects", err) {
		for _, s := range list {
			subjects[s.ID] = s.Name
		}
	}
	users := map[int]string{}
	if list, err := client.Users(ctx); check("users", err) {
		for _, u := range list {
			users[u.ID] = u.FirstName + " " + u.LastName
		}
	}

	section("Me")
	if me, err := client.Me(ctx); check("me", err) {
		fmt.Printf("%s %s (%s)\n", me.User.FirstName, me.User.LastName, me.Account.Login)
	}
	if class, err := client.Class(ctx); check("class", err) {
		fmt.Printf("class %d%s, school year ends %s\n", class.Number, class.Symbol, class.EndSchoolYear)
	}
	if classrooms, err := client.Classrooms(ctx); check("classrooms", err) {
		fmt.Printf("%d classrooms in the school\n", len(classrooms))
	}

	section("Grades")
	if grades, err := client.Grades(ctx); check("grades", err) {
		categories := map[int]librus.GradeCategory{}
		if list, err := client.GradeCategories(ctx); check("grade categories", err) {
			for _, c := range list {
				categories[c.ID] = c
			}
		}
		for _, g := range grades {
			category := categories[g.CategoryID]
			fmt.Printf("%-22s %-3s %s (%s, weight %d) by %s\n", subjects[g.SubjectID], g.Grade,
				g.Date.Format("02.01"), category.Name, category.Weight, users[g.AddedByID])
			// Comments cost one request per comment, so they are fetched only when a grade has any.
			comments, err := g.Comments(ctx)
			if check("grade comments", err) {
				for _, c := range comments {
					fmt.Printf("    comment: %s\n", c.Text)
				}
			}
		}
	}

	section("Point grades")
	if grades, err := client.PointGrades(ctx); check("point grades", err) {
		for _, g := range grades {
			category, err := g.Category(ctx)
			if !check("point grade category", err) {
				continue
			}
			fmt.Printf("%-22s %s/%d %s (%s)\n", subjects[g.SubjectID], g.Grade, category.MaxPoints,
				g.Date.Format("02.01"), category.Name)
		}
	}

	section("Averages")
	if averages, err := client.Averages(ctx); check("averages", err) {
		for _, a := range averages {
			fmt.Printf("%-22s semester 1: %.2f (%d counted)\n", subjects[a.SubjectID], a.Semester1.Value, a.Semester1.Count)
		}
	}

	section("Timetable this week")
	if lessons, err := client.Timetable(ctx, time.Now()); check("timetable", err) {
		var day time.Time
		for _, l := range lessons {
			if !l.Date.Equal(day) {
				day = l.Date
				fmt.Println(day.Format("Monday 02.01"))
			}
			status := ""
			switch {
			case l.IsCanceled:
				status = " [canceled]"
			case l.IsSubstitution:
				status = " [substitution]"
				// The replaced lesson's subject is fetched only for substitutions.
				if l.Original != nil {
					if s, err := l.Original.Subject(ctx); err == nil {
						status += " instead of " + s.Name
					}
				}
			}
			fmt.Printf("  %d. %s-%s %s%s\n", l.LessonNo, l.Start.Format("15:04"), l.End.Format("15:04"), l.Subject.Name, status)
		}
	}

	section("Timetable changes")
	if list, err := client.Substitutions(ctx); check("substitutions", err) {
		for _, s := range list {
			if s.IsCanceled {
				fmt.Printf("%s lesson %d canceled: %s\n", s.OrgDate.Format("02.01"), s.OrgLessonNo, subjects[s.OrgSubjectID])
			} else {
				fmt.Printf("%s lesson %d: %s with %s instead of %s with %s\n", s.Date.Format("02.01"), s.LessonNo,
					subjects[s.SubjectID], users[s.TeacherID], subjects[s.OrgSubjectID], users[s.OrgTeacherID])
			}
		}
	}

	section("Upcoming tests, trips and other events")
	if events, err := client.Events(ctx); check("events", err) {
		categories := map[int]string{}
		if list, err := client.EventCategories(ctx); check("event categories", err) {
			for _, c := range list {
				categories[c.ID] = c.Name
			}
		}
		sort.Slice(events, func(i, j int) bool { return events[i].Date.Before(events[j].Date) })
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		for _, e := range events {
			if e.Date.Before(today) {
				continue
			}
			fmt.Printf("%s %-24s %-20s %s\n", e.Date.Format("02.01"), categories[e.CategoryID], subjects[e.SubjectID], firstLine(e.Content))
		}
	}

	section("Homework")
	if list, err := client.HomeworkAssignments(ctx); check("homework", err) {
		for _, h := range list {
			fmt.Printf("due %s from %s: %s\n", h.DueDate.Format("02.01"), users[h.TeacherID], h.Topic)
		}
	}

	section("Attendance")
	if attendances, err := client.Attendances(ctx); check("attendances", err) {
		types := map[int]string{}
		if list, err := client.AttendanceTypes(ctx); check("attendance types", err) {
			for _, t := range list {
				types[t.ID] = t.Name
			}
		}
		counts := map[string]int{}
		for _, a := range attendances {
			counts[types[a.TypeID]]++
		}
		for name, n := range counts {
			fmt.Printf("%-20s %d\n", name, n)
		}
	}
	if list, err := client.Justifications(ctx); check("justifications", err) {
		for _, j := range list {
			fmt.Printf("justification %s - %s (%s), %d absences\n", j.DateFrom.Format("02.01"), j.DateTo.Format("02.01"), j.Status, j.JustifiedAbsences)
		}
	}

	section("Next days off")
	if days, err := client.SchoolFreeDays(ctx); check("school free days", err) {
		shown := 0
		for _, d := range days {
			if d.DateTo.Before(time.Now()) || shown == 3 {
				continue
			}
			shown++
			fmt.Printf("%s - %s %s\n", d.DateFrom.Format("02.01"), d.DateTo.Format("02.01"), strings.Join(strings.Fields(d.Name), " "))
		}
	}
}

func section(title string) {
	fmt.Printf("\n== %s ==\n", title)
}

// check reports whether err is nil and prints it otherwise.
func check(what string, err error) bool {
	if err == nil {
		return true
	}
	var apiErr *librus.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("(%s unavailable: %s)\n", what, apiErr.Message)
	} else {
		fmt.Printf("(%s failed: %v)\n", what, err)
	}
	return false
}

func firstLine(s string) string {
	line, _, _ := strings.Cut(s, "\n")
	if len(line) > 60 {
		line = line[:57] + "..."
	}
	return line
}
