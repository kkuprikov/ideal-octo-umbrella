package streamprocessor

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func GenerateScreenshot(ctx context.Context, url string, dir string, c chan string) {
	sliced := strings.Split(url, "/")
	streamID := sliced[len(sliced)-1]

	var pathParts = []string{dir, "/", streamID, ".jpg"}
	path := strings.Join(pathParts, "")

	var args = []string{"-y", "-nostdin", "-skip_frame", "nokey", "-i", url,
		"-vsync", "0", "-r", "30", "-f", "image2", "-update", "1", path}

	cmd := exec.Command("ffmpeg", args...)
	err := cmd.Start()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Screenshot task for %s started", url)

	select {
	case control := <-c:
		controls := strings.Split(control, "_")
		if (controls[0] == "snapshots") && (controls[1] == url) {
			if err := cmd.Process.Kill(); err != nil {
				log.Fatal("failed to kill screenshot task: ", err)
			}
		}
		return
	case <-ctx.Done():
		fmt.Println("ctx.Done() in GenerateScreenshot")

		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill screenshot task: ", err)
		}
		return
	}
}
