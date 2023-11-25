package mi

import "time"

type StartInfo struct {
	StartID     string     `json:"start_id"`
	TaskID      string     `json:"task_id"`
	UpdatedTime time.Time  `json:"updated_time"`
	Start       *time.Time `json:"start"`
}
