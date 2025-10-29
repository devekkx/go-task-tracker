package storage

import (
	"github.com/devekkx/go-task-tracker/internal/models"
)

// Store holds all application data.
type Store struct {
	Tasks     []models.Task     `json:"tasks"`
	TodoLists []models.TodoList `json:"todo_lists"`
	path      string
}
