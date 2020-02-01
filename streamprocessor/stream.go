package streamprocessor

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

func GetStreamData(ctx context.Context, url string, control chan string, out chan string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

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

	go StoreData(ctx, out)

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		select {
		case c := <-control:
			controls := strings.Split(c, "_")
			if (controls[0] == "stats") && (controls[1] == url) {
				if err := cmd.Process.Kill(); err != nil {
					log.Fatal("failed to kill stats task: ", err)
				}
				fmt.Println("Stats task for %s stopped", url)
				close(out)
				return
			}
		case <-ctx.Done():
			fmt.Println("ctx.Done() in GetStreamData")

			if err := cmd.Process.Kill(); err != nil {
				log.Fatal("failed to kill stats task: ", err)
			}
			fmt.Println("Stats task for %s stopped", url)
			close(out)
			return
		default:
			out <- scanner.Text() + "," + strconv.FormatInt(time.Now().UTC().Unix(), 10)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading ffmpeg output:", err)
	}
}
