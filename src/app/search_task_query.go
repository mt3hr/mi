package mi

type SearchTaskQuery struct {
	Board      string
	Tags       []string
	Word       string
	CheckState CheckState
	SortType   SortType
}
