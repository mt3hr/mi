package mi

import "time"

type BoardInfo struct {
	BoardInfoID string    `json:"board_info_id"`
	TaskID      string    `json:"task_id"`
	UpdatedTime time.Time `json:"updated_time"`
	BoardName   string    `json:"board_name"`
}
