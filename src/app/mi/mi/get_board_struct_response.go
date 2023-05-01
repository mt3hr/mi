package miapp

type GetBoardStructResponse struct {
	Errors      []string    `json:"errors"`
	BoardStruct interface{} `json:"board_struct"`
}
