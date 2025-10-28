package models

import "time"

// TodoItem is a single item in a TodoList.
type TodoItem struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TodoList groups related TodoItems.
type TodoList struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Items     []TodoItem `json:"items"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
