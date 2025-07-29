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

	userInfo, err := librusApi.WithReauth(
		func() (*librusApi.UserInfo, error) {
			return librusApi.GetUserInfo(client)
		},
		func() error {
			return client.Login(login, password)
		},
		1)

	fmt.Println(userInfo)
}
