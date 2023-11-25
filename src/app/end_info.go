package mi

import "time"

type EndInfo struct {
	EndID       string     `json:"end_id"`
	TaskID      string     `json:"task_id"`
	UpdatedTime time.Time  `json:"updated_time"`
	End         *time.Time `json:"end"`
}
