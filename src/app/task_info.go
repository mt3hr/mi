package mi

type TaskInfo struct {
	Task           *Task           `json:"task"`
	TaskTitleInfo  *TaskTitleInfo  `json:"task_title_info"`
	CheckStateInfo *CheckStateInfo `json:"check_state_info"`
	LimitInfo      *LimitInfo      `json:"limit_info"`
	BoardInfo      *BoardInfo      `json:"board_info"`
}
