package miapp

import mi "github.com/mt3hr/mi/src/app"

type UpdateTaskRequest struct {
	TaskInfo *mi.TaskInfo `json:"task_info"`
}
