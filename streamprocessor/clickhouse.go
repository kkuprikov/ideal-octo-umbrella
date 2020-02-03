package streamprocessor

import (
	"context"
	"fmt"
	"sync"
)

func storeData(ctx context.Context, out chan string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	// read from channel, form some chunks and send to database
	fmt.Println("Time to store some data!")
	for {
		select {
		case v, ok := <-out:
			if !ok {
				fmt.Println("Stat channel closed, exiting...")
				return
			}
			fmt.Println("Received ", v) // store data
		case <-ctx.Done():
			fmt.Println("ctx.Done() in StoreData")
			return
		}
	}
}
