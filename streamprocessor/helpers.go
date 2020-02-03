package streamprocessor

import (
	"fmt"
	"os/exec"
	"strings"
)

// type for API params

type streamParams struct {
	Stream string `json:"stream"`
}

// types for unmarshalling nested JSON from ingester

type stream struct {
	Name    string `json:"name"`
	Publish struct {
		Active bool `json:"active"`
	} `json:"publish"`
}

type ingesterStreams struct {
	Code    int      `json:"code"`
	Server  int      `json:"server"`
	Streams []stream `json:"streams"`
}

func GetID(url string) string {
	sliced := strings.Split(url, "/")
	streamID := sliced[len(sliced)-1]

	if strings.Contains(streamID, "?") {
		streamID = strings.Split(streamID, "?")[0]
	}
	return streamID
}

func killProcess(cmd *exec.Cmd) {
	if err := cmd.Process.Kill(); err != nil {
		fmt.Println("Failed to kill screenshot task: ", err)
	}
}
