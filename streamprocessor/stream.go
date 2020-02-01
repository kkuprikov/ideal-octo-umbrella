package streamprocessor

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

func GetStreamData(ctx context.Context, url string, control chan string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	streamID := GetID(url)

	args := []string{"-loglevel", "error", "-select_streams", "v:0", "-show_frames",
		"-show_entries", "frame=key_frame,pkt_duration_time,pkt_size,height,repeat_pict", "-of", "csv", url}

	cmd := exec.Command("ffprobe", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	out := make(chan string)

	go StoreData(ctx, out, wg)
	go readData(stdout, out, streamID)

	select {
	case <-ctx.Done():
		stdout.Close()
		fmt.Println("ctx done in GetStreamData!")
	}
}

func readData(stdout io.ReadCloser, out chan string, streamID string) {
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		out <- scanner.Text() + "," + strconv.FormatInt(time.Now().UTC().Unix(), 10) + "," + streamID
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading ffmpeg output:", err)
	}
}
