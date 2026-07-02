package models

import (
	"fmt"
	"time"
)

// Priority represents task urgency level.
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Status represents the lifecycle state of a task.
type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
)

// Task is the core domain entity.
type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Priority    Priority   `json:"priority"`
	Status      Status     `json:"status"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// NewTask creates a Task with defaults.
func NewTask(title, description string, priority Priority) *Task {
	now := time.Now()
	return &Task{
		ID:          generateID("task"),
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      StatusPending,
		Tags:        []string{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// MarkDone sets the task status to done.
func (t *Task) MarkDone() {
	t.Status = StatusDone
	t.UpdatedAt = time.Now()
}

// MarkInProgress sets the task status to in-progress.
func (t *Task) MarkInProgress() {
	t.Status = StatusInProgress
	t.UpdatedAt = time.Now()
}

// MarkPending resets the task to pending.
func (t *Task) MarkPending() {
	t.Status = StatusPending
	t.UpdatedAt = time.Now()
}

// SetDueDate assigns a due date to the task.
func (t *Task) SetDueDate(d time.Time) {
	t.DueDate = &d
	t.UpdatedAt = time.Now()
}

// AddTag appends tag to the task if it is not already present.
func (t *Task) AddTag(tag string) {
	for _, existing := range t.Tags {
		if existing == tag {
			return
		}
	}
	t.Tags = append(t.Tags, tag)
	t.UpdatedAt = time.Now()
}

// IsOverdue returns true when the task has a past due date and is not done.
func (t *Task) IsOverdue() bool {
	if t.DueDate == nil {
		return false
	}
	return time.Now().After(*t.DueDate) && t.Status != StatusDone
}

// DaysUntilDue returns the number of days until the due date, or -1 if unset.
func (t *Task) DaysUntilDue() int {
	if t.DueDate == nil {
		return -1
	}
	return int(time.Until(*t.DueDate).Hours() / 24)
}

// Validate checks that the task has required fields
func (t *Task) Validate() error {
	if t.Title == "" {
		return fmt.Errorf("task title cannot be empty")
	}
	switch t.Priority {
	case PriorityLow, PriorityMedium, PriorityHigh:
	default:
		return fmt.Errorf("invalid priority %q: must be low, medium, or high", t.Priority)
	}
	switch t.Status {
	case StatusPending, StatusInProgress, StatusDone:
	default:
		return fmt.Errorf("invalid status %q", t.Status)
	}
	return nil
}

// ValidPriority parses and validates a priority string.
func ValidPriority(s string) (Priority, error) {
	switch Priority(s) {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return Priority(s), nil
	default:
		return "", fmt.Errorf("invalid priority %q: must be low, medium, or high", s)
	}
}

// ValidStatus parses and validates a status string.
func ValidStatus(s string) (Status, error) {
	switch Status(s) {
	case StatusPending, StatusInProgress, StatusDone:
		return Status(s), nil
	default:
		return "", fmt.Errorf("invalid status %q: must be pending, in-progress, or done", s)
	}
}

// HasTag reports whether the task has the given tag.
func (t *Task) HasTag(tag string) bool {
	for _, existing := range t.Tags {
		if existing == tag {
			return true
		}
	}
	return false
}
