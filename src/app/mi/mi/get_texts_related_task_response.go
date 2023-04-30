package miapp

import (
	text "github.com/mt3hr/rykv/text"
)

type GetTextsRelatedTaskResponse struct {
	Errors []string     `json:"errors"`
	Texts  []*text.Text `json:"texts"`
}
