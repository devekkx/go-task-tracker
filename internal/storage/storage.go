package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
		return fmt.Errorf("failed to parse data file %s: %w", s.path, err)
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

// DataPath returns the absolute path to the backing JSON file.
// Useful for diagnostics and export workflows.
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
// All string fields are optional; empty string means no filter on that dimension.
type FilterOptions struct {
	Status   string
	Priority string
	Tag      string
	Search   string
	Archived *bool // nil = exclude archived, true = only archived, false = all
	SortBy   string // title, priority, due, created
}

// ListTasks returns tasks matching the given filter options.
func (s *Store) ListTasks(opts FilterOptions) []models.Task {
	result := make([]models.Task, 0, len(s.Tasks))
	for _, t := range s.Tasks {
		if matchesFilter(t, opts) {
			result = append(result, t)
		}
	}
	sortTasks(result, opts.SortBy)
	return result
}

func sortTasks(tasks []models.Task, by string) {
	switch by {
	case "title":
		sort.Slice(tasks, func(i, j int) bool {
			return strings.ToLower(tasks[i].Title) < strings.ToLower(tasks[j].Title)
		})
	case "priority":
		order := map[models.Priority]int{models.PriorityHigh: 0, models.PriorityMedium: 1, models.PriorityLow: 2}
		sort.Slice(tasks, func(i, j int) bool {
			return order[tasks[i].Priority] < order[tasks[j].Priority]
		})
	case "due":
		sort.Slice(tasks, func(i, j int) bool {
			if tasks[i].DueDate == nil {
				return false
			}
			if tasks[j].DueDate == nil {
				return true
			}
			return tasks[i].DueDate.Before(*tasks[j].DueDate)
		})
	case "created":
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
		})
	}
}

func matchesFilter(t models.Task, opts FilterOptions) bool {
	if opts.Archived == nil {
		if t.Archived {
			return false
		}
	} else if *opts.Archived != t.Archived {
		return false
	}
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
	if opts.Tag != "" { // exact case-sensitive tag match
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
	ArchivedTasks   int
	TotalTodoLists  int
	TotalTodoItems  int
	DoneTodoItems    int
	PendingTodoItems int
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
		if t.Archived {
			stats.ArchivedTasks++
		}
	}
	for _, l := range s.TodoLists {
		stats.TotalTodoItems += l.TotalItems()
		stats.DoneTodoItems += l.DoneItems()
		stats.PendingTodoItems += l.PendingItems()
	}
	return stats
}

// ArchiveTask marks the given task as archived.
func (s *Store) ArchiveTask(id string) error {
	task, err := s.GetTask(id)
	if err != nil {
		return err
	}
	task.Archive()
	return s.UpdateTask(task)
}

// UnarchiveTask restores a previously archived task.
func (s *Store) UnarchiveTask(id string) error {
	task, err := s.GetTask(id)
	if err != nil {
		return err
	}
	task.Unarchive()
	return s.UpdateTask(task)
}

// ClearDoneTasks removes all tasks with status done and returns the count removed.
func (s *Store) ClearDoneTasks() (int, error) {
	remaining := make([]models.Task, 0, len(s.Tasks))
	removed := 0
	for _, t := range s.Tasks {
		if t.Status == models.StatusDone {
			removed++
		} else {
			remaining = append(remaining, t)
		}
	}
	if removed == 0 {
		return 0, nil
	}
	s.Tasks = remaining
	return removed, s.save()
}

// ExportJSON writes the store contents to the given file path as JSON.
func (s *Store) ExportJSON(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}
	return nil
}

// ImportJSON reads tasks and todo lists from a JSON export file,
// merging them into the current store (skipping duplicate IDs).
func (s *Store) ImportJSON(path string) (int, int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read import file: %w", err)
	}
	var src Store
	if err := json.Unmarshal(data, &src); err != nil {
		return 0, 0, fmt.Errorf("failed to parse import file: %w", err)
	}
	importedTasks := 0
	for i := range src.Tasks {
		if err := src.Tasks[i].Validate(); err != nil {
			return 0, 0, fmt.Errorf("invalid task in import: %w", err)
		}
		duplicate := false
		for _, existing := range s.Tasks {
			if existing.ID == src.Tasks[i].ID {
				duplicate = true
				break
			}
		}
		if !duplicate {
			s.Tasks = append(s.Tasks, src.Tasks[i])
			importedTasks++
		}
	}
	importedLists := 0
	for i := range src.TodoLists {
		if err := src.TodoLists[i].Validate(); err != nil {
			return 0, 0, fmt.Errorf("invalid todo list in import: %w", err)
		}
		s.TodoLists = append(s.TodoLists, src.TodoLists[i])
		importedLists++
	}
	return importedTasks, importedLists, s.save()
}

// SearchTodoLists returns todo lists whose name contains query (case-insensitive).
func (s *Store) SearchTodoLists(query string) []models.TodoList {
	q := strings.ToLower(query)
	result := make([]models.TodoList, 0)
	for _, l := range s.TodoLists {
		if strings.Contains(strings.ToLower(l.Name), q) {
			result = append(result, l)
		}
	}
	return result
}

// BulkMarkDone marks all tasks matching opts as done and returns the count updated.
func (s *Store) BulkMarkDone(opts FilterOptions) (int, error) {
	tasks := s.ListTasks(opts)
	if len(tasks) == 0 {
		return 0, nil
	}
	for _, t := range tasks {
		task := t
		task.MarkDone()
		if err := s.UpdateTask(&task); err != nil {
			return 0, err
		}
	}
	return len(tasks), nil
}

// CopyTask creates a duplicate of the given task with a new ID.
func (s *Store) CopyTask(id string) (*models.Task, error) {
	src, err := s.GetTask(id)
	if err != nil {
		return nil, err
	}
	copied := models.NewTask(src.Title+" (copy)", src.Description, src.Priority)
	copied.Tags = append([]string{}, src.Tags...)
	if src.DueDate != nil {
		d := *src.DueDate
		copied.SetDueDate(d)
	}
	if err := s.AddTask(copied); err != nil {
		return nil, err
	}
	return copied, nil
}

// TaskCount returns the total number of tasks including archived.
func (s *Store) TaskCount() int { return len(s.Tasks) }

// TodoListCount returns the total number of todo lists.
func (s *Store) TodoListCount() int { return len(s.TodoLists) }
