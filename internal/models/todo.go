package models

import (
	"fmt"
	"time"
)

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

// AddItem appends a new item to the list and returns it.
func (l *TodoList) AddItem(content string) TodoItem {
	item := NewTodoItem(content)
	l.Items = append(l.Items, item)
	l.UpdatedAt = time.Now()
	return item
}

// CheckItem marks the item with itemID as done.
func (l *TodoList) CheckItem(itemID string) error {
	for i := range l.Items {
		if l.Items[i].ID == itemID {
			l.Items[i].Done = true
			l.Items[i].UpdatedAt = time.Now()
			l.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("item %q not found in list %q", itemID, l.Name)
}

// UncheckItem marks the item with itemID as not done.
func (l *TodoList) UncheckItem(itemID string) error {
	for i := range l.Items {
		if l.Items[i].ID == itemID {
			l.Items[i].Done = false
			l.Items[i].UpdatedAt = time.Now()
			l.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("item %q not found in list %q", itemID, l.Name)
}

// RemoveItem deletes the item with itemID from the list.
func (l *TodoList) RemoveItem(itemID string) error {
	for i, item := range l.Items {
		if item.ID == itemID {
			l.Items = append(l.Items[:i], l.Items[i+1:]...)
			l.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("item %q not found in list %q", itemID, l.Name)
}

// GetItem returns a pointer to the item with itemID.
func (l *TodoList) GetItem(itemID string) (*TodoItem, error) {
	for i := range l.Items {
		if l.Items[i].ID == itemID {
			return &l.Items[i], nil
		}
	}
	return nil, fmt.Errorf("item %q not found", itemID)
}

// TotalItems returns the number of items in the list.
func (l *TodoList) TotalItems() int { return len(l.Items) }

// DoneItems returns the number of completed items.
func (l *TodoList) DoneItems() int {
	count := 0
	for _, item := range l.Items {
		if item.Done {
			count++
		}
	}
	return count
}

// PendingItems returns the number of incomplete items.
func (l *TodoList) PendingItems() int { return l.TotalItems() - l.DoneItems() }

// Progress returns the percentage of completed items.
func (l *TodoList) Progress() float64 {
	if l.TotalItems() == 0 {
		return 0
	}
	return float64(l.DoneItems()) / float64(l.TotalItems()) * 100
}

// Validate verifies the todo list has required fields
func (l *TodoList) Validate() error {
	if l.Name == "" {
		return fmt.Errorf("todo list name cannot be empty")
	}
	return nil
}
