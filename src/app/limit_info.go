package mi

import "time"

type LimitInfo struct {
	LimitID     string     `json:"limit_id"`
	TaskID      string     `json:"task_id"`
	UpdatedTime time.Time  `json:"updated_time"`
	Limit       *time.Time `json:"limit"`
}
