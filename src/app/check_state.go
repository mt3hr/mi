package mi

type CheckState int

const (
	NoCheckOnly CheckState = 0
	CheckOnly   CheckState = 1
	All         CheckState = 2
)
