package miapp

type GetBoardNamesResponse struct {
	Errors     []string `json:"errors"`
	BoardNames []string `json:"board_names"`
}
