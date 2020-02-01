package streamprocessor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
)

type ControlMessage struct {
	Host    string
	Type    string // stats or snapshots
	Message string // start or stop
	Url     string
}

func Subscribe(ctx context.Context, control chan string, stats chan string, wg *sync.WaitGroup) {
	redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		fmt.Println("redis db cannot be set up:", err)
		return
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDb,
	})

	fmt.Println("subscribe started")
	psNewMessage := redisClient.Subscribe(os.Getenv("REDIS_CHANNEL"))

	for {
		msg, err := psNewMessage.ReceiveMessage()
		if err != nil {
			fmt.Println("error on redis subscription: %s", err)
		}
		go processMessage(ctx, msg.Payload, control, stats, wg)
	}

}

func processMessage(ctx context.Context, msg string, control chan string, stats chan string, wg *sync.WaitGroup) {
	var data ControlMessage

	if err := json.Unmarshal([]byte(msg), &data); err != nil {
		fmt.Println("JSON unmarshalling error: ", err)
		return
	}

	fmt.Println("received control msg: %s", data)

	// start or stop task, depending on message type
	switch data.Message {
	case "stop":
		control <- data.Type + "_" + data.Url
	case "start":
		switch data.Type {
		case "stats":
			go GetStreamData(ctx, data.Url, control, stats, wg)
		case "snapshots":
			go GenerateScreenshot(ctx, data.Url, os.Getenv("SNAPSHOT_DIR"), control, wg)
		}
	}
}
