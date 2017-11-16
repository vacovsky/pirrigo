package pirri

import "encoding/json"

func reactToGpioMessage(m []byte) {
	task := Task{}
	json.Unmarshal(m, &task)
	task.execute()
}
