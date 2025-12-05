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

// save flushes the store to disk
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

// DataPath returns the absolute path to the backing JSON file.
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
// All fields are optional; empty string means no filter.
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
	if opts.Search != "" {
		q := strings.ToLower(opts.Search)
		if !strings.Contains(strings.ToLower(t.Title), q) && !strings.Contains(strings.ToLower(t.Description), q) {
			return false
		}
	}
	if opts.Status != "" && string(t.Status) != opts.Status {
		return false
	}
	if opts.Priority != "" && string(t.Priority) != opts.Priority {
		return false
	}
	if opts.Tag != "" { // case-sensitive tag match
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

// --- Todo CRUD ---

// AddTodoList persists a new TodoList.
func (s *Store) AddTodoList(list *models.TodoList) error {
	if err := list.Validate(); err != nil {
		return err
	}
	s.TodoLists = append(s.TodoLists, *list)
	return s.save()
}

// GetTodoList returns the list with the given ID.
func (s *Store) GetTodoList(id string) (*models.TodoList, error) {
	for i := range s.TodoLists {
		if s.TodoLists[i].ID == id {
			return &s.TodoLists[i], nil
		}
	}
	return nil, fmt.Errorf("todo list %q not found", id)
}

// UpdateTodoList replaces the stored list with the given one.
func (s *Store) UpdateTodoList(list *models.TodoList) error {
	for i := range s.TodoLists {
		if s.TodoLists[i].ID == list.ID {
			s.TodoLists[i] = *list
			return s.save()
		}
	}
	return fmt.Errorf("todo list %q not found", list.ID)
}

// DeleteTodoList removes the list with the given ID.
func (s *Store) DeleteTodoList(id string) error {
	for i, l := range s.TodoLists {
		if l.ID == id {
			s.TodoLists = append(s.TodoLists[:i], s.TodoLists[i+1:]...)
			return s.save()
		}
	}
	return fmt.Errorf("todo list %q not found", id)
}

// ListTodoLists returns a copy of all todo lists.
func (s *Store) ListTodoLists() []models.TodoList {
	result := make([]models.TodoList, len(s.TodoLists))
	copy(result, s.TodoLists)
	return result
}

// Stats summarises task and todo counts.
type Stats struct {
	TotalTasks      int
	PendingTasks    int
	InProgressTasks int
	DoneTasks       int
	OverdueTasks    int
	TotalTodoLists  int
	TotalTodoItems  int
	DoneTodoItems   int
}

// GetStats computes aggregate statistics across all tasks and todo lists.
func (s *Store) GetStats() Stats {
	stats := Stats{
		TotalTasks:     len(s.Tasks),
		TotalTodoLists: len(s.TodoLists),
	}
	for _, t := range s.Tasks {
		switch t.Status {
		case models.StatusPending:
			stats.PendingTasks++
		case models.StatusInProgress:
			stats.InProgressTasks++
		case models.StatusDone:
			stats.DoneTasks++
		}
		if t.IsOverdue() {
			stats.OverdueTasks++
		}
	}
	for _, l := range s.TodoLists {
		stats.TotalTodoItems += l.TotalItems()
		stats.DoneTodoItems += l.DoneItems()
	}
	return stats
}
