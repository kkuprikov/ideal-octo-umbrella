package main

import (
	"context"
	"log"
	"streamprocessor"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	control := make(chan string)
	out := make(chan string)

	streamprocessor.Subscribe(ctx, control, out)
}
