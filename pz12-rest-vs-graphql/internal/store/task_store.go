package store

type Task struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Done        bool    `json:"done"`
}

var Tasks = []*Task{
	{
		ID:          "t_001",
		Title:       "First task",
		Description: strPtr("REST API example"),
		Done:        false,
	},
	{
		ID:          "t_002",
		Title:       "Second task",
		Description: strPtr("GraphQL example"),
		Done:        true,
	},
}

func strPtr(s string) *string {
	return &s
}