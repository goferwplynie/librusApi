// Prints this week's timetable, marking canceled lessons and substitutions.
package main

import (
	"context"
	"fmt"
	"os"
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

	lessons, err := client.Timetable(ctx, time.Now())
	if err != nil {
		panic(err)
	}

	var day time.Time
	for _, l := range lessons {
		if !l.Date.Equal(day) {
			day = l.Date
			fmt.Printf("\n%s\n", day.Format("Monday 02.01"))
		}

		status := ""
		switch {
		case l.IsCanceled:
			status = " [canceled]"
		case l.IsSubstitution:
			status = " [substitution]"
		}

		fmt.Printf("  %d. %s-%s  %-30s %s %s%s\n",
			l.LessonNo, l.Start.Format("15:04"), l.End.Format("15:04"),
			l.Subject.Name, l.Teacher.FirstName, l.Teacher.LastName, status)
		if l.Note != "" {
			fmt.Printf("     note: %s\n", l.Note)
		}
	}
}
