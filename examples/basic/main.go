package main

import (
	"fmt"
	"os"

	librus "github.com/goferwplynie/librusApi"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	login := os.Getenv("LIBRUS_LOGIN")
	password := os.Getenv("LIBRUS_PASSWORD")

	client := librus.New()
	if err := client.Login(login, password); err != nil {
		panic(err)
	}

	me, err := librus.WithReauth(client.Me, func() error {
		return client.Login(login, password)
	}, 1)
	if err != nil {
		panic(err)
	}

	fmt.Println("Hello", me.User.FirstName)
}
