package main

import (
	"context"
	"fmt"
	"os"

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

	me, err := client.Me(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Hello", me.User.FirstName)
}
