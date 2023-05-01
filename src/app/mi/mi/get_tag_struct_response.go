package miapp

type GetTagStructResponse struct {
	Errors    []string    `json:"errors"`
	TagStruct interface{} `json:"tag_struct"`
}
