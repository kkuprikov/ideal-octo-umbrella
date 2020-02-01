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
	"strings"
	"sync"
	"time"
)

func GetStreamData(ctx context.Context, url string, control chan string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	defer fmt.Println("Stats task for %s stopped", url)

	streamID := GetID(url)

	args := []string{"-loglevel", "error", "-select_streams", "v:0", "-show_frames",
		"-show_entries", "frame=key_frame,pkt_duration_time,pkt_size,height,repeat_pict", "-of", "csv", url}

	cmd := exec.Command("ffprobe", args...)

	stdout, err := cmd.StdoutPipe()
	defer stdout.Close()

	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	out := make(chan string)

	go StoreData(ctx, out, wg)
	go readData(stdout, out, streamID, wg)

	select {
	case <-ctx.Done():
		fmt.Println("ctx done in GetStreamData!")
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill stats task: ", err)
		}
		return
	case c := <-control:
		controls := strings.Split(c, "_")
		if (controls[0] == "stats") && (controls[1] == url) {
			if err := cmd.Process.Kill(); err != nil {
				log.Fatal("failed to kill stats task: ", err)
			}
			return
		}
	}
}

func readData(stdout io.ReadCloser, out chan string, streamID string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	defer close(out)
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		out <- scanner.Text() + "," + strconv.FormatInt(time.Now().UTC().Unix(), 10) + "," + streamID
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading ffmpeg output:", err)
	}
}
