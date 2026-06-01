package pirri

import (
	"time"
)

func ListenForTasks() {
	defer WG.Done()
	for {
		ORQMutex.Lock()
		var task *Task
		if len(OfflineRunQueue) > 0 {
			task = OfflineRunQueue[0]
			OfflineRunQueue = OfflineRunQueue[1:]
		}
		ORQMutex.Unlock()

		if task != nil {
			task.execute()
		}
		time.Sleep(time.Duration(1000) * time.Millisecond)
	}
}
