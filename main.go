package main

import (
	"fmt"
	"os"

	"log"

	"github.com/goferwplynie/librusApi/client"
	"github.com/goferwplynie/librusApi/users"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	login := os.Getenv("LIBRUS_LOGIN")
	password := os.Getenv("LIBRUS_PASSWORD")

	client := client.New()

	client.Login(login, password)

	grades, err := users.GetUsers(client)
	if err != nil {
		panic(err)
	}
	for _, v := range grades {
		fmt.Println(v)
	}
}
