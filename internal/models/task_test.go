package models_test

import (
	"testing"
	"time"

	"github.com/devekkx/go-task-tracker/internal/models"
)

func TestNewTask_defaults(t *testing.T) {
	task := models.NewTask("Write tests", "", models.PriorityMedium)
	if task.Title != "Write tests" {
		t.Errorf("expected title 'Write tests', got %q", task.Title)
	}
	if task.Status != models.StatusPending {
		t.Errorf("expected status pending, got %q", task.Status)
	}
	if task.Priority != models.PriorityMedium {
		t.Errorf("expected priority medium, got %q", task.Priority)
	}
	if task.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestTask_MarkDone(t *testing.T) {
	task := models.NewTask("Test", "", models.PriorityLow)
	task.MarkDone()
	if task.Status != models.StatusDone {
		t.Errorf("expected done, got %q", task.Status)
	}
}

func TestTask_IsOverdue(t *testing.T) {
	task := models.NewTask("Old task", "", models.PriorityHigh)
	past := time.Now().Add(-24 * time.Hour)
	task.SetDueDate(past)
	if !task.IsOverdue() {
		t.Error("expected task to be overdue")
	}
	task.MarkDone()
	if task.IsOverdue() {
		t.Error("done task should not be overdue")
	}
}

func TestTask_AddTag(t *testing.T) {
	task := models.NewTask("Tagged", "", models.PriorityLow)
	task.AddTag("work")
	task.AddTag("work") // duplicate
	if len(task.Tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(task.Tags))
	}
}

func TestTask_Validate(t *testing.T) {
	task := models.NewTask("", "", models.PriorityLow)
	if err := task.Validate(); err == nil {
		t.Error("expected error for empty title")
	}
}

func TestValidPriority(t *testing.T) {
	for _, valid := range []string{"low", "medium", "high"} {
		if _, err := models.ValidPriority(valid); err != nil {
			t.Errorf("expected %q to be valid: %v", valid, err)
		}
	}
	if _, err := models.ValidPriority("urgent"); err == nil {
		t.Error("expected error for invalid priority")
	}
}

func TestTask_DaysUntilDue_unset(t *testing.T) {
	task := models.NewTask("No due", "", models.PriorityLow)
	if task.DaysUntilDue() != -1 {
		t.Error("expected -1 for unset due date")
	}
}

func TestTask_HasTag(t *testing.T) {
	task := models.NewTask("Tag test", "", models.PriorityLow)
	task.AddTag("work")
	if !task.HasTag("work") {
		t.Error("expected HasTag to return true for existing tag")
	}
	if task.HasTag("personal") {
		t.Error("expected HasTag to return false for missing tag")
	}
}

func TestTask_RemoveTag(t *testing.T) {
	task := models.NewTask("Remove tag", "", models.PriorityLow)
	task.AddTag("work")
	task.AddTag("urgent")
	task.RemoveTag("work")
	if task.HasTag("work") {
		t.Error("expected work tag to be removed")
	}
	if !task.HasTag("urgent") {
		t.Error("expected urgent tag to remain")
	}
}

func TestTask_ArchiveUnarchive(t *testing.T) {
	task := models.NewTask("Archive me", "", models.PriorityMedium)
	if task.Archived {
		t.Error("new task should not be archived")
	}
	task.Archive()
	if !task.Archived {
		t.Error("expected task to be archived")
	}
	task.Unarchive()
	if task.Archived {
		t.Error("expected task to be unarchived")
	}
}

func TestTask_Validate_whitespaceTitle(t *testing.T) {
	task := models.NewTask("   ", "", models.PriorityLow)
	if err := task.Validate(); err == nil {
		t.Error("expected error for whitespace-only title")
	}
}

func TestTask_TagsCount(t *testing.T) {
	task := models.NewTask("Tags", "", models.PriorityLow)
	if task.TagsCount() != 0 {
		t.Errorf("expected 0 tags, got %d", task.TagsCount())
	}
	task.AddTag("a")
	task.AddTag("b")
	if task.TagsCount() != 2 {
		t.Errorf("expected 2 tags, got %d", task.TagsCount())
	}
}

func TestTask_MarkInProgress(t *testing.T) {
	task := models.NewTask("Start me", "", models.PriorityMedium)
	task.MarkInProgress()
	if task.Status != models.StatusInProgress {
		t.Errorf("expected in-progress, got %q", task.Status)
	}
}

func TestTask_MarkPending(t *testing.T) {
	task := models.NewTask("Reset me", "", models.PriorityLow)
	task.MarkDone()
	task.MarkPending()
	if task.Status != models.StatusPending {
		t.Errorf("expected pending, got %q", task.Status)
	}
}

func TestValidStatus(t *testing.T) {
	for _, valid := range []string{"pending", "in-progress", "done"} {
		if _, err := models.ValidStatus(valid); err != nil {
			t.Errorf("expected %q to be valid: %v", valid, err)
		}
	}
	if _, err := models.ValidStatus("cancelled"); err == nil {
		t.Error("expected error for invalid status")
	}
}
