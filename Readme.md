# Go Librus API

Librus API library. I wanted to make some widgets for my system but most of the libraries are outdated.Also i couldn't find any documentation for this API so this Library can help a lot with apps that want to use it

<details>
<summary>Implemented Endpoints Checklist</summary>

- [x] `Users` (teachers)
- [x] `Grades`  
- [ ] `Me`  
- [ ] `PointGrades`  
- [ ] `Subjects`  
- [ ] `Lessons`  
- [ ] `Realizations` (maybe) 
- [ ] `Timetables` (badly designed so maybe in future) 
- [ ] `TimetableEntries` 
- [ ] `Classes`  
- [ ] `Classrooms`  
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

## Usage

basic example for getting all available users
```go
package main

import (
	"fmt"
	"os"

	"github.com/goferwplynie/librusApi"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	login := os.Getenv("LIBRUS_LOGIN")
	password := os.Getenv("LIBRUS_PASSWORD")

	client := librusApi.New()

	client.Login(login, password)

	users, err := librusApi.WithReauth(
		func() ([]*librusApi.User, error) {
			return librusApi.GetUsers(client)
		},
		func() error {
			return client.Login(login, password)
		},
		1)

	for _, v := range users {
		fmt.Println(v)
	}
}
```
first you have to create new client and use login method with your login and password.
```go
	client := librusApi.New()

	client.Login(login, password)
```
after this you should be able to get all available resources. 1 session is 10-15 minutes long so i recommend using `librusApi.withReauth(method, loginFunc, retries)`.
```go
	users, err := librusApi.WithReauth(
		func() ([]*users.User, error) {
			return librusApi.GetUsers(client)
		},
		func() error {
			return client.Login(login, password)
		},
		1)
```
All functions that are supposed to get resources from api take `*librusApi.Client` as argument.
```go
librusApi.GetUsers(client)
```

Librus API doesn't explictly say if session expired so when program gets `Access Denied` it assumes that session expired and throws special error `SessionExpiredError`.
If you decide to not use `librusApi.withReauth` you will have to take care of handling this error case because if you are not logged in client won't make any requests.
