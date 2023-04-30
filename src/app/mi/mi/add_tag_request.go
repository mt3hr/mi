package miapp

type AddTagRequest struct {
	TaskID string `json:"task_id"`
	Tag    string `json:"tag"`
}
