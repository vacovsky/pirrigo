package pirri

import (
	"time"
)

func ListenForTasks() {
	defer WG.Done()
	for {
		ORQMutex.Lock()
		q := OfflineRunQueue
		ORQMutex.Unlock()

		var task *Task
		if len(q) > 0 {
			ORQMutex.Lock()
			task, OfflineRunQueue = OfflineRunQueue[len(OfflineRunQueue)-1],
				OfflineRunQueue[:len(OfflineRunQueue)-1]
			ORQMutex.Unlock()
			task.execute()
		}
		time.Sleep(time.Duration(1000) * time.Millisecond)
	}
}
