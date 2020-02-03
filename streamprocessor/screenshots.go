package streamprocessor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) screenshotsStart(ctx context.Context, wg *sync.WaitGroup, tasks *SafeTasks) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		var rec streamParams
		err := json.NewDecoder(r.Body).Decode(&rec)

		if err != nil || rec.Stream == "" {
			w.WriteHeader(400)
			return
		}

		id := rec.Stream

		NewScreenshotTask(ctx, id, tasks, wg)

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "0")

	}
}

func NewScreenshotTask(ctx context.Context, id string, tasks *SafeTasks, wg *sync.WaitGroup) {
	url := os.Getenv("INGESTER_URL") + "/" + id
	dir := os.Getenv("SNAPSHOT_DIR")

	if !screenshotsTaskExists(tasks, id) {
		ch := addScreenshotsTask(tasks, id)
		go generateScreenshot(ctx, url, dir, ch, wg)
	}
}

func (s *Server) screenshotsStop(ctx context.Context, wg *sync.WaitGroup, tasks *SafeTasks) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		var rec streamParams
		err := json.NewDecoder(r.Body).Decode(&rec)

		if err != nil || rec.Stream == "" {
			w.WriteHeader(400)
			return
		}

		id := rec.Stream

		if screenshotsTaskExists(tasks, id) {
			removeScreenshotsTask(tasks, id)
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "0")
	}
}

func generateScreenshot(ctx context.Context, url string, dir string, c chan string, wg *sync.WaitGroup) {
	// wg.Add(1)
	// defer wg.Done()

	streamID := GetID(url)

	var pathParts = []string{dir, "/", streamID, ".jpg"}
	path := strings.Join(pathParts, "")

	var args = []string{"-y", "-nostdin", "-skip_frame", "nokey", "-i", url,
		"-vsync", "0", "-r", "30", "-f", "image2", "-update", "1", path}

	cmd := exec.Command("ffmpeg", args...)

	if err := cmd.Start(); err != nil {
		fmt.Printf("Can't start screenshot task %s\n", url)
		return
	}

	go waitForFinish(ctx, url, cmd, c, wg)
	cmd.Wait()
	c <- "done"
	fmt.Printf("Screenshot task for %s finished with no signal\n", url)
	return
}

func waitForFinish(ctx context.Context, url string, cmd *exec.Cmd, c chan string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	fmt.Printf("Screenshot task for %s started, state: %s\n", url, cmd.ProcessState)

	select {
	case <-c:
		killProcess(cmd)
		fmt.Printf("Screenshot task for %s stopped by channel\n", url)
		return
	case <-ctx.Done():
		fmt.Printf("ctx.Done() in GenerateScreenshot\n")
		killProcess(cmd)
		fmt.Printf("Screenshot task for %s stopped by signal\n", url)
		return
	}
}
