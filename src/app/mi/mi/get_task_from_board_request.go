package miapp

import (
	mi "github.com/mt3hr/mi/src/app"
)

type GetTaskFromBoardRequest struct {
	Query *mi.SearchTaskQuery `json:"query"`
}
