package miapp

type GetBoardStructResponse struct {
	Errors      []string `json:"errors"`
	BoardStruct string   `json:"board_struct"`
}
