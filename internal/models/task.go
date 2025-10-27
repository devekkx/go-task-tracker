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
