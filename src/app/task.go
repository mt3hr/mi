package mi

import "time"

type Task struct {
	TaskID      string    `json:"task_id"`
	CreatedTime time.Time `json:"created_time"`
}
