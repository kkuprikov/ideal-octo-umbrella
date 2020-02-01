package streamprocessor

import (
	"context"
	"fmt"
)

func StoreData(ctx context.Context, out chan string) {
	// read from channel, form some chunks and send to database
	fmt.Println("Time to store some data!")
}
