package streamprocessor

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func GetRunningStreams(url string) []string {
	var s ingesterStreams
	m := make([]string, 0, 50)

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &s)

	if err != nil {
		log.Fatal(err)
	}

	for _, ss := range s.Streams {
		if ss.Publish.Active == true {
			m = append(m, ss.Name)
		}
	}
	return m
}
