package streamprocessor

import "context"

func StoreData(ctx context.Context, out chan string) {
	// read from channel, form some chunks and send to database
}
