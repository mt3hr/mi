package mi

import "time"

type CheckStateInfo struct {
	CheckStateID string    `json:"check_state_id"`
	TaskID       string    `json:"task_id"`
	UpdatedTime  time.Time `json:"updated_time"`
	IsChecked    bool      `json:"is_checked"`
}
