package miapp

import (
	tag "github.com/mt3hr/rykv/tag"
)

type GetTagsRelatedTaskResponse struct {
	Errors []string   `json:"errors"`
	Tags   []*tag.Tag `json:"tags"`
}
