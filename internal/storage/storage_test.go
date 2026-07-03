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

func TestStore_UnarchiveTask(t *testing.T) {
	s := tempStore(t)
	task := models.NewTask("Restore me", "", models.PriorityMedium)
	_ = s.AddTask(task)
	_ = s.ArchiveTask(task.ID)

	if err := s.UnarchiveTask(task.ID); err != nil {
		t.Fatalf("UnarchiveTask: %v", err)
	}

	tasks := s.ListTasks(storage.FilterOptions{})
	found := false
	for _, tsk := range tasks {
		if tsk.ID == task.ID {
			found = true
		}
	}
	if !found {
		t.Error("unarchived task should appear in default list")
	}
}

func TestStore_ListTasks_sortByTitle(t *testing.T) {
	s := tempStore(t)
	_ = s.AddTask(models.NewTask("Zebra task", "", models.PriorityLow))
	_ = s.AddTask(models.NewTask("Apple task", "", models.PriorityLow))
	_ = s.AddTask(models.NewTask("Mango task", "", models.PriorityLow))

	tasks := s.ListTasks(storage.FilterOptions{SortBy: "title"})
	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}
	if tasks[0].Title != "Apple task" {
		t.Errorf("expected first task to be Apple task, got %q", tasks[0].Title)
	}
}

func TestStore_ListTasks_sortByPriority(t *testing.T) {
	s := tempStore(t)
	_ = s.AddTask(models.NewTask("Low task", "", models.PriorityLow))
	_ = s.AddTask(models.NewTask("High task", "", models.PriorityHigh))
	_ = s.AddTask(models.NewTask("Med task", "", models.PriorityMedium))

	tasks := s.ListTasks(storage.FilterOptions{SortBy: "priority"})
	if tasks[0].Title != "High task" {
		t.Errorf("expected first task to be High task, got %q", tasks[0].Title)
	}
}

func TestStore_ClearDoneTasks(t *testing.T) {
	s := tempStore(t)
	t1 := models.NewTask("Pending", "", models.PriorityLow)
	t2 := models.NewTask("Done", "", models.PriorityLow)
	t2.MarkDone()
	_ = s.AddTask(t1)
	_ = s.AddTask(t2)

	n, err := s.ClearDoneTasks()
	if err != nil {
		t.Fatalf("ClearDoneTasks: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 removed, got %d", n)
	}
	tasks := s.ListTasks(storage.FilterOptions{})
	if len(tasks) != 1 {
		t.Errorf("expected 1 remaining task, got %d", len(tasks))
	}
}

func TestStore_ExportImportRoundtrip(t *testing.T) {
	s := tempStore(t)
	_ = s.AddTask(models.NewTask("Export me", "desc", models.PriorityHigh))
	_ = s.AddTodoList(models.NewTodoList("My List"))

	exportPath := t.TempDir() + "/export.json"
	if err := s.ExportJSON(exportPath); err != nil {
		t.Fatalf("ExportJSON: %v", err)
	}

	s2 := tempStore(t)
	tasks, lists, err := s2.ImportJSON(exportPath)
	if err != nil {
		t.Fatalf("ImportJSON: %v", err)
	}
	if tasks != 1 {
		t.Errorf("expected 1 imported task, got %d", tasks)
	}
	if lists != 1 {
		t.Errorf("expected 1 imported list, got %d", lists)
	}
}

func TestStore_SearchTodoLists(t *testing.T) {
	s := tempStore(t)
	_ = s.AddTodoList(models.NewTodoList("Work tasks"))
	_ = s.AddTodoList(models.NewTodoList("Personal tasks"))
	_ = s.AddTodoList(models.NewTodoList("Grocery list"))

	results := s.SearchTodoLists("tasks")
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}
