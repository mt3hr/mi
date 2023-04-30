package miapp

import text "github.com/mt3hr/rykv/text"

type GetTextResponse struct {
	Errors []string   `json:"errors"`
	Text   *text.Text `json:"text"`
}
