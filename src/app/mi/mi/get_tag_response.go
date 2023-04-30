package miapp

import tag "github.com/mt3hr/rykv/tag"

type GetTagResponse struct {
	Errors []string `json:"errors"`
	Tag    *tag.Tag `json:"tag"`
}
