package mi

type SearchTaskQuery struct {
	Board      string     `json:"board"`
	Tags       []string   `json:"tags"`
	Word       string     `json:"word"`
	CheckState CheckState `json:"check_state"`
	SortType   SortType   `json:"sort_type"`
}
