package miapp

import mi "github.com/mt3hr/mi/src/app"

type GetTaskFromBoardResponse struct {
	Errors      []string       `json:"errors"`
	BoardsTasks []*mi.TaskInfo `json:"boards_tasks"`
}
