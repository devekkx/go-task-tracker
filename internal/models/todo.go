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

// NewTodoList creates a TodoList with the given name.
func NewTodoList(name string) *TodoList {
	now := time.Now()
	return &TodoList{
		ID:        generateID("todo"),
		Name:      name,
		Items:     []TodoItem{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewTodoItem creates a TodoItem with the given content.
func NewTodoItem(content string) TodoItem {
	now := time.Now()
	return TodoItem{
		ID:        generateID("item"),
		Content:   content,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
