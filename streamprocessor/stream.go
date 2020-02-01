package streamprocessor

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func GetStreamData(ctx context.Context, url string, control chan string, out chan string) {
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

				close(out)
				return
			}
		case <-ctx.Done():
			fmt.Println("ctx.Done() in GetStreamData")

			if err := cmd.Process.Kill(); err != nil {
				log.Fatal("failed to kill stats task: ", err)
			}

			close(out)
			return
		default:
			out <- scanner.Text() + "," + time.Now().UTC().Unix()
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading ffmpeg output:", err)
	}
}
