package streamprocessor

import (
	"context"
	"fmt"
	"sync"
)

func StoreData(ctx context.Context, out chan string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	// read from channel, form some chunks and send to database
	fmt.Println("Time to store some data!")
	for {
		select {
		case v, ok := <-out:
			fmt.Println("Received ", v)
			if !ok {
				fmt.Println("Stat channel closed, exiting...")
				return
			}
		case <-ctx.Done():
			fmt.Println("ctx.Done() in StoreData")
			return
		}
	}
}
