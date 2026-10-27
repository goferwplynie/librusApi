# Go Librus API

Go client for [Librus Synergia](https://synergia.librus.pl): grades, timetable, attendance and more for a student account.

Librus has no public API documentation, and most existing libraries are outdated or scrape HTML. This one uses the JSON API that Synergia exposes for the web session.

## Installation

```sh
go get github.com/goferwplynie/librusApi
```

The package name is `librus`:

```go
import librus "github.com/goferwplynie/librusApi"
```

## Usage

```go
ctx := context.Background()

client := librus.New()
if err := client.Login(ctx, login, password); err != nil {
	log.Fatal(err)
}

grades, err := client.Grades(ctx)
if err != nil {
	log.Fatal(err)
}
for _, g := range grades {
	fmt.Println(g.Grade, g.Date.Format(time.DateOnly))
}
```

Runnable examples are in [`examples`](examples): [`basic`](examples/basic/main.go) logs in and greets you, [`timetable`](examples/timetable/main.go) prints this week's timetable and [`all`](examples/all/main.go) shows everything the library can fetch. They read `LIBRUS_LOGIN` and `LIBRUS_PASSWORD` from a `.env` file in the repo root:

```sh
go run ./examples/basic
go run ./examples/timetable
go run ./examples/all
```

### Available methods

| Method | Returns |
|---|---|
| `Me(ctx)` | logged in account |
| `Grades(ctx)` | grades (`Grade.Teacher(ctx)`, `Grade.Subject(ctx)`, `Grade.Category(ctx)` and `Grade.Comments(ctx)` fetch what the grade points to) |
| `GradeCategories(ctx)`, `GradeCategory(ctx, id)` | grade categories with weights |
| `PointGrades(ctx)`, `PointGrade(ctx, id)` | point grades (`Category`, `Comments`, `Subject` and `Teacher` work like on `Grade`) |
| `PointGradeCategories(ctx)`, `PointGradeComments(ctx)` | point grade categories (with point ranges) and comments |
| `GradeComments(ctx)`, `GradeComment(ctx, id)` | teachers' comments to grades |
| `Timetable(ctx, day)` | lessons of the week (Mon-Sun) containing `day`, incl. canceled lessons and substitutions (`OriginalLesson.Subject(ctx)` fetches the subject of a replaced lesson) |
| `TimetableDay(ctx, day)` | lessons of a single day |
| `Events(ctx)`, `Event(ctx, id)` | calendar entries ("terminarz"): tests, quizzes, trips... (`Category`, `Subject` and `CreatedBy` fetch what the event points to) |
| `EventCategories(ctx)`, `EventCategory(ctx, id)` | kinds of calendar events, e.g. "Sprawdzian" |
| `HomeworkAssignments(ctx)`, `HomeworkAssignment(ctx, id)` | homework set by teachers |
| `Attendances(ctx)` | attendance at lessons (`Attendance.Type(ctx)`, `Attendance.Lesson(ctx)`) |
| `AttendanceTypes(ctx)`, `AttendanceType(ctx, id)` | kinds of attendance (presence, absence, lateness...) |
| `Lessons(ctx)`, `Lesson(ctx, id)` | lessons of the class: which teacher teaches which subject |
| `Justifications(ctx)`, `Justification(ctx, id)` | absence justifications sent by parents |
| `SchoolFreeDays(ctx)` | school days off, with `Includes(day)` |
| `Substitutions(ctx)`, `Substitution(ctx, id)` | timetable changes: canceled, moved and substituted lessons |
| `Averages(ctx)` | averages per subject, computed locally (see `CalculateAverages` for the rules) |
| `Users(ctx)`, `User(ctx, id)` | teachers and other users |
| `Subjects(ctx)`, `Subject(ctx, id)` | school subjects (`Grade.Subject(ctx)` fetches the subject of a grade) |
| `Class(ctx)` | student's class |
| `Classrooms(ctx)` | school classrooms |
| `School(ctx)` | school info |
| `Get(ctx, resource)` | raw `*http.Response` for endpoints not wrapped yet, e.g. `Get(ctx, "Attendances")` |

### Sessions

Librus sessions expire after ~10-15 minutes. The client keeps the credentials passed to `Login` in memory and logs in again automatically when the session expires. A `Client` is safe for concurrent use.

### Errors

- `ErrLoginFailed` - wrong login or password.
- `ErrNotLoggedIn` - a method was called before `Login`.
- `*APIError` - Librus responded with status >= 400. Has `StatusCode`, `Code` and `Message`.

```go
var apiErr *librus.APIError
if errors.As(err, &apiErr) && apiErr.StatusCode == 403 {
	// feature disabled by the school, e.g. Code == "LuckyNumberIsNotActive"
}
```

## Roadmap

### Done (1.0)
- [x] Timetable
- [x] `Me`, `Users`, `Classes`, `Classrooms`, `Schools`
- [x] `Subjects`, `Grades`
- [x] Grade categories and comments
- [x] Point grades
- [x] Averages - computed locally from grade weights, because `Grades/Averages` can be disabled by the school
- [x] Attendance and lessons
- [x] Justifications, free days and substitutions
- [x] Tests and quizzes (terminarz) and homework

### Later
- [ ] Combined calendar (timetable, tests, trips, free days)
- [ ] `SchoolNotices`, `SchoolTrips`, `ParentTeacherConferences`, `ClassFreeDays`, `Notes`
- [ ] `TextGrades`, `DescriptiveGrades`, `BehaviourGrades`
- [ ] `LuckyNumbers`
- [ ] Messages - `Messages` returns `Insufficient scopes`
