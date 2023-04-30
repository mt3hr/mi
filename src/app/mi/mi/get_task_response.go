package miapp

import mi "github.com/mt3hr/mi/src/app"

type GetTaskResponse struct {
	Errors   []string     `json:"errors"`
	TaskInfo *mi.TaskInfo `json:"task_info"`
}
