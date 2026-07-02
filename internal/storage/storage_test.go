package storage_test

import (
	"os"
	"testing"

	"github.com/devekkx/go-task-tracker/internal/models"
	"github.com/devekkx/go-task-tracker/internal/storage"
)

func tempStore(t *testing.T) *storage.Store {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	s, err := storage.New()
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return s
}

func TestStore_AddGetTask(t *testing.T) {
	s := tempStore(t)
	task := models.NewTask("Test", "desc", models.PriorityHigh)
	if err := s.AddTask(task); err != nil {
		t.Fatalf("AddTask: %v", err)
	}
	got, err := s.GetTask(task.ID)
	if err != nil {
		t.Fatalf("GetTask: %v", err)
	}
	if got.Title != "Test" {
		t.Errorf("expected title 'Test', got %q", got.Title)
	}
}

func TestStore_DeleteTask(t *testing.T) {
	s := tempStore(t)
	task := models.NewTask("Del", "", models.PriorityLow)
	_ = s.AddTask(task)
	if err := s.DeleteTask(task.ID); err != nil {
		t.Fatalf("DeleteTask: %v", err)
	}
	if _, err := s.GetTask(task.ID); err == nil {
		t.Error("expected error for deleted task")
	}
}

func TestStore_ListTasks_filter(t *testing.T) {
	s := tempStore(t)
	_ = s.AddTask(models.NewTask("High task", "", models.PriorityHigh))
	_ = s.AddTask(models.NewTask("Low task", "", models.PriorityLow))
	tasks := s.ListTasks(storage.FilterOptions{Priority: "high"})
	if len(tasks) != 1 {
		t.Errorf("expected 1 high task, got %d", len(tasks))
	}
}

func TestStore_AddGetTodoList(t *testing.T) {
	s := tempStore(t)
	list := models.NewTodoList("Shopping")
	if err := s.AddTodoList(list); err != nil {
		t.Fatalf("AddTodoList: %v", err)
	}
	got, err := s.GetTodoList(list.ID)
	if err != nil {
		t.Fatalf("GetTodoList: %v", err)
	}
	if got.Name != "Shopping" {
		t.Errorf("expected 'Shopping', got %q", got.Name)
	}
}

func TestStore_Persistence(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	s1, _ := storage.New()
	task := models.NewTask("Persist me", "", models.PriorityMedium)
	_ = s1.AddTask(task)

	s2, _ := storage.New()
	got, err := s2.GetTask(task.ID)
	if err != nil {
		t.Fatalf("task not persisted: %v", err)
	}
	if got.Title != "Persist me" {
		t.Errorf("expected 'Persist me', got %q", got.Title)
	}
	_ = os.Unsetenv("HOME")
}

func TestStore_UpdateTask(t *testing.T) {
	s := tempStore(t)
	task := models.NewTask("Original", "", models.PriorityLow)
	_ = s.AddTask(task)
	task.Title = "Updated"
	if err := s.UpdateTask(task); err != nil {
		t.Fatalf("UpdateTask: %v", err)
	}
	got, _ := s.GetTask(task.ID)
	if got.Title != "Updated" {
		t.Errorf("expected 'Updated', got %q", got.Title)
	}
}

func TestStore_ArchiveTask(t *testing.T) {
	s := tempStore(t)
	task := models.NewTask("Archive me", "", models.PriorityLow)
	_ = s.AddTask(task)

	if err := s.ArchiveTask(task.ID); err != nil {
		t.Fatalf("ArchiveTask: %v", err)
	}

	// archived task should not appear in default list
	tasks := s.ListTasks(storage.FilterOptions{})
	for _, tsk := range tasks {
		if tsk.ID == task.ID {
			t.Error("archived task should not appear in default list")
		}
	}

	// should appear when explicitly requesting archived
	archived := true
	archivedTasks := s.ListTasks(storage.FilterOptions{Archived: &archived})
	if len(archivedTasks) != 1 {
		t.Errorf("expected 1 archived task, got %d", len(archivedTasks))
	}
}
