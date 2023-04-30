package mi

import "time"

type TaskTitleInfo struct {
	TaskTitleID string    `json:"task_title_id"`
	TaskID      string    `json:"task_id"`
	UpdatedTime time.Time `json:"updated_time"`
	Title       string    `json:"title"`
}
