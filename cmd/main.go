package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kkuprikov/streamprocessor-go/streamprocessor"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	control := make(chan string)

	var wg sync.WaitGroup

	go streamprocessor.Subscribe(ctx, control, &wg)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-termChan
	fmt.Println("Shutdown signal received")
	cancelFunc() // Signal cancellation to context.Context
	wg.Wait()

	fmt.Println("All workers done, shutting down!")
}
