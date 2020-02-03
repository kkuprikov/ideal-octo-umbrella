package streamprocessor

import (
	"fmt"
	"sync"
)

type SafeTasks struct {
	Stats       map[string](chan string)
	Screenshots map[string](chan string)
	Mux         sync.Mutex
}

func addStatsTask(tasks *SafeTasks, id string) chan string {
	tasks.Mux.Lock()
	defer tasks.Mux.Unlock()

	tasks.Stats[id] = make(chan string)
	return tasks.Stats[id]
}

func addScreenshotsTask(tasks *SafeTasks, id string) chan string {
	tasks.Mux.Lock()
	defer tasks.Mux.Unlock()
	fmt.Println("task add", id, tasks.Screenshots[id])

	tasks.Screenshots[id] = make(chan string)
	return tasks.Screenshots[id]
}

func removeStatsTask(tasks *SafeTasks, id string) {
	tasks.Mux.Lock()
	defer tasks.Mux.Unlock()
	tasks.Stats[id] <- "done"
	close(tasks.Stats[id])
	delete(tasks.Stats, id)
	return
}

func removeScreenshotsTask(tasks *SafeTasks, id string) {
	tasks.Mux.Lock()
	defer tasks.Mux.Unlock()
	fmt.Println("task remove", id, tasks.Screenshots[id])
	tasks.Screenshots[id] <- "done"
	close(tasks.Screenshots[id])
	delete(tasks.Screenshots, id)
	return
}

func statsTaskExists(tasks *SafeTasks, id string) bool {
	tasks.Mux.Lock()
	defer tasks.Mux.Unlock()

	return tasks.Stats[id] != nil
}

func screenshotsTaskExists(tasks *SafeTasks, id string) bool {
	tasks.Mux.Lock()
	defer tasks.Mux.Unlock()
	fmt.Println("task exists", id, tasks.Screenshots[id])
	return tasks.Screenshots[id] != nil
}
