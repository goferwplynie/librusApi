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

A runnable example is in [`examples/basic`](examples/basic/main.go). It reads `LIBRUS_LOGIN` and `LIBRUS_PASSWORD` from a `.env` file in the repo root:

```sh
go run ./examples/basic
```

### Available methods

| Method | Returns |
|---|---|
| `Me(ctx)` | logged in account |
| `Grades(ctx)` | grades (`Grade.Teacher(ctx)` fetches who added it) |
| `Users(ctx)`, `User(ctx, id)` | teachers and other users |
| `Class(ctx)` | student's class |
| `Classrooms(ctx)` | school classrooms |
| `School(ctx)` | school info |
| `Get(ctx, resource)` | raw `*http.Response` for endpoints not wrapped yet, e.g. `Get(ctx, "Subjects")` |

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

### Done
- [x] `Me`, `Users`, `Classes`, `Classrooms`, `Schools`
- [x] `Grades` (IDs only, no subject or category names yet)

### 1.0
- [ ] Timetable - `Timetables?weekStart=YYYY-MM-DD`
- [ ] Subjects - `Subjects`
- [ ] Grade categories and comments - `Grades/Categories`, `Grades/Comments`
- [ ] Point grades - `PointGrades`
- [ ] Averages - computed locally from grade weights, because `Grades/Averages` can be disabled by the school
- [ ] Tests and quizzes (terminarz) - `HomeWorks`, `HomeWorks/Categories`
- [ ] Homework - `HomeWorkAssignments`
- [ ] Attendance - `Attendances`, `Attendances/Types`, `Lessons`
- [ ] Justifications - `Justifications`
- [ ] Free days - `SchoolFreeDays`
- [ ] Substitutions - `Calendars/Substitutions`

### Later
- [ ] Combined calendar (timetable, tests, trips, free days)
- [ ] `SchoolNotices`, `SchoolTrips`, `ParentTeacherConferences`, `ClassFreeDays`, `Notes`
- [ ] `TextGrades`, `DescriptiveGrades`, `BehaviourGrades`
- [ ] `LuckyNumbers`
- [ ] Messages - `Messages` returns `Insufficient scopes`
