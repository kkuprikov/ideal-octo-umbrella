package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/julienschmidt/httprouter"
	"github.com/kkuprikov/streamprocessor-go/streamprocessor"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	appPort := os.Getenv("PORT")

	if appPort == "" {
		log.Fatal("Application port not specified")
	}

	ingesterUrl := os.Getenv("INGESTER_API_URL")

	if ingesterUrl == "" {
		log.Fatal("Ingester URL not specified")
	}

	var (
		wg    sync.WaitGroup
		tasks streamprocessor.SafeTasks
	)

	tasks.Stats = make(map[string](chan string))
	tasks.Screenshots = make(map[string](chan string))

	for _, id := range streamprocessor.GetRunningStreams(ingesterUrl) {
		streamprocessor.NewScreenshotTask(ctx, id, &tasks, &wg)
		streamprocessor.NewStatsTask(ctx, id, &tasks, &wg)
	}

	s := &streamprocessor.Server{}
	s.Router = httprouter.New()

	go s.Start(ctx, appPort, &wg, &tasks)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-termChan
	fmt.Println("Shutdown signal received")
	cancelFunc() // Signal cancellation to context.Context
	wg.Wait()

	fmt.Println("All workers done, shutting down!")
}
