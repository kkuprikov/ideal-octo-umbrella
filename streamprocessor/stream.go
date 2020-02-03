package streamprocessor

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) statsStart(ctx context.Context, wg *sync.WaitGroup, tasks *SafeTasks) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		var rec streamParams
		err := json.NewDecoder(r.Body).Decode(&rec)

		if err != nil || rec.Stream == "" {
			w.WriteHeader(400)
			return
		}

		id := rec.Stream

		url := os.Getenv("INGESTER_URL") + "/" + id
		NewStatsTask(ctx, url, tasks, wg)

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "0")
	}
}

func NewStatsTask(ctx context.Context, id string, tasks *SafeTasks, wg *sync.WaitGroup) {
	if !statsTaskExists(tasks, id) {
		// fmt.Printf("NEW STATS %s\n", id)
		ch := addStatsTask(tasks, id)
		url := os.Getenv("INGESTER_URL") + "/" + id
		go getStreamData(ctx, url, ch, wg)
	}
}

func (s *Server) statsStop(ctx context.Context, wg *sync.WaitGroup, tasks *SafeTasks) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		var rec streamParams
		err := json.NewDecoder(r.Body).Decode(&rec)

		if err != nil || rec.Stream == "" {
			w.WriteHeader(400)
			return
		}

		id := rec.Stream

		if statsTaskExists(tasks, id) {
			removeStatsTask(tasks, id)
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "0")
	}
}

func getStreamData(ctx context.Context, url string, control chan string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	streamID := GetID(url)

	args := []string{"-loglevel", "error", "-select_streams", "v:0", "-show_frames",
		"-show_entries", "frame=key_frame,pkt_duration_time,pkt_size,height,repeat_pict", "-of", "csv", url}

	cmd := exec.Command("ffprobe", args...)

	stdout, err := cmd.StdoutPipe()
	defer stdout.Close()

	if err != nil {
		fmt.Printf("Failed to set stdout for stats task %s", url)
		return
	}

	if err = cmd.Start(); err != nil {
		fmt.Printf("Failed to start stats task %s", url)
		return
	}
	out := make(chan string)

	go readData(stdout, out, streamID, wg)
	go storeData(ctx, out, wg)

	select {
	case <-ctx.Done():
		fmt.Println("ctx done in GetStreamData!")
		killProcess(cmd)
		fmt.Printf("Stats task for %s stopped by context\n", url)
		return
	case <-control:
		killProcess(cmd)
		fmt.Printf("Stats task for %s stopped by channel\n", url)
		return
	}
}

func readData(stdout io.ReadCloser, out chan string, streamID string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	defer close(out)
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		out <- strconv.FormatInt(time.Now().UTC().Unix(), 10) + "," + streamID + "," + scanner.Text()
	}
	fmt.Println("STOPPED")
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading ffmpeg output:", err)
	}
}
