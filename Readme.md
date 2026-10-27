# Go Librus API

Librus API library. I wanted to make some widgets for my system but most of the libraries are outdated.Also i couldn't find any documentation for this API so this Library can help a lot with apps that want to use it

<details>
<summary>Implemented Endpoints Checklist</summary>

- [x] `Users` (teachers)
- [x] `Grades`  
- [x] `Me`  
- - [ ] `GradeAverage`
- [ ] `PointGrades`  
- [ ] `Subjects`  
- [ ] `Lessons`  
- [ ] `Realizations` (maybe) 
- [ ] `Timetables` (badly designed so maybe in future) 
- [ ] `TimetableEntries` 
- [x] `Classes`  
- [x] `Classrooms`  
- [x] `Schools`
- [ ] `ClassFreeDays`  
- [ ] `SchoolFreeDays`  
- [ ] `Attendances`  
- [ ] `Justifications`  
- [ ] `HomeWorks`  
- [ ] `HomeWorkAssignments`  
- [ ] `SchoolTrips`  
- [ ] `Messages`  
- [ ] `ParentTeacherConferences`  
- [ ] `ExamResult`  
- [ ] `Substitutions`  
- [ ] `Calendars`  

</details>

There are a lot of other data endpoints but there is a problem with access to it or how it looks.

## Installation

```sh
go get github.com/goferwplynie/librusApi
```

The package is called `librus`.

## Usage

basic example for getting all available users (see also [`examples/basic`](examples/basic/main.go))
```go
package main

import (
	"fmt"
	"os"

	librus "github.com/goferwplynie/librusApi"
)

func main() {
	login := os.Getenv("LIBRUS_LOGIN")
	password := os.Getenv("LIBRUS_PASSWORD")

	client := librus.New()
	if err := client.Login(login, password); err != nil {
		panic(err)
	}

	users, err := librus.WithReauth(client.Users, func() error {
		return client.Login(login, password)
	}, 1)
	if err != nil {
		panic(err)
	}

	for _, v := range users {
		fmt.Println(v.FirstName, v.LastName)
	}
}
```
first you have to create new client and log in with your login and password.
```go
	client := librus.New()
	if err := client.Login(login, password); err != nil {
		panic(err)
	}
```
after this you can get resources using client methods:
```go
me, err := client.Me()
grades, err := client.Grades()
teacher, err := grades[0].Teacher()
```
For endpoints that aren't wrapped yet you can use `client.Get("Subjects")`, which returns the raw `*http.Response`.

1 session is 10-15 minutes long so i recommend using `librus.WithReauth(method, loginFunc, retries)`. It takes any method without arguments (like `client.Users`), so for methods with arguments wrap them in a closure:
```go
	user, err := librus.WithReauth(func() (*librus.User, error) {
		return client.User(id)
	}, func() error {
		return client.Login(login, password)
	}, 1)
```

Librus API doesn't explictly say if session expired so when program gets `Access Denied` it assumes that session expired and returns special error `SessionExpiredError` (check it with `librus.IsSessionExpired(err)`).
If you decide to not use `librus.WithReauth` you will have to take care of handling this error case because if you are not logged in client won't make any requests.
