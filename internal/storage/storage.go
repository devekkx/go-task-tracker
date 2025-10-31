package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devekkx/go-task-tracker/internal/models"
)

const dataFileName = "data.json"

// Store holds all application data.
type Store struct {
	Tasks     []models.Task     `json:"tasks"`
	TodoLists []models.TodoList `json:"todo_lists"`
	path      string
}

// New creates a Store backed by the user's home directory.
func New() (*Store, error) {
	dir, err := dataDir()
	if err != nil {
		return nil, fmt.Errorf("failed to determine data directory: %w", err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}
	s := &Store{
		Tasks:     []models.Task{},
		TodoLists: []models.TodoList{},
		path:      filepath.Join(dir, dataFileName),
	}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func dataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".task-tracker"), nil
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to read data file: %w", err)
	}
	if err := json.Unmarshal(data, s); err != nil {
		return fmt.Errorf("failed to parse data file: %w", err)
	}
	return nil
}

// save writes the store to disk
func (s *Store) save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize data: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0644); err != nil {
		return fmt.Errorf("failed to write data file: %w", err)
	}
	return nil
}

// DataPath returns the path to the backing JSON file.
func (s *Store) DataPath() string { return s.path }

// AddTask persists a new task.
func (s *Store) AddTask(task *models.Task) error {
	if err := task.Validate(); err != nil {
		return err
	}
	s.Tasks = append(s.Tasks, *task)
	return s.save()
}

// GetTask returns the task with the given ID.
func (s *Store) GetTask(id string) (*models.Task, error) {
	for i := range s.Tasks {
		if s.Tasks[i].ID == id {
			return &s.Tasks[i], nil
		}
	}
	return nil, fmt.Errorf("task %q not found", id)
}

// UpdateTask replaces the stored task with the given one.
func (s *Store) UpdateTask(task *models.Task) error {
	if err := task.Validate(); err != nil {
		return err
	}
	for i := range s.Tasks {
		if s.Tasks[i].ID == task.ID {
			s.Tasks[i] = *task
			return s.save()
		}
	}
	return fmt.Errorf("task %q not found", task.ID)
}

// DeleteTask removes the task with the given ID.
func (s *Store) DeleteTask(id string) error {
	for i, t := range s.Tasks {
		if t.ID == id {
			s.Tasks = append(s.Tasks[:i], s.Tasks[i+1:]...)
			return s.save()
		}
	}
	return fmt.Errorf("task %q not found", id)
}

// FilterOptions controls which tasks ListTasks returns.
type FilterOptions struct {
	Status   string
	Priority string
	Tag      string
	Search   string
}

// ListTasks returns tasks matching the given filter options.
func (s *Store) ListTasks(opts FilterOptions) []models.Task {
	result := make([]models.Task, 0)
	for _, t := range s.Tasks {
		if matchesFilter(t, opts) {
			result = append(result, t)
		}
	}
	return result
}

func matchesFilter(t models.Task, opts FilterOptions) bool {
	if opts.Status != "" && string(t.Status) != opts.Status {
		return false
	}
	if opts.Priority != "" && string(t.Priority) != opts.Priority {
		return false
	}
	if opts.Tag != "" {
		found := false
		for _, tag := range t.Tags {
			if tag == opts.Tag {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
