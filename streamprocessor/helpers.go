package streamprocessor

import "strings"

func GetID(url string) string {
	sliced := strings.Split(url, "/")
	streamID := sliced[len(sliced)-1]

	if strings.Contains(streamID, "?") {
		streamID = strings.Split(streamID, "?")[0]
	}
	return streamID
}
