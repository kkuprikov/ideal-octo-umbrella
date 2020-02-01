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
	for v := range out {
		fmt.Println("Received ", v)
	}
}
