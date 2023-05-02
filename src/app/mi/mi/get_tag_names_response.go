package miapp

type GetTagNamesResponse struct {
	Errors   []string `json:"errors"`
	TagNames []string `json:"tag_names"`
}
