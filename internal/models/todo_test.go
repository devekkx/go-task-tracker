package models_test

import (
	"testing"

	"github.com/devekkx/go-task-tracker/internal/models"
)

func TestNewTodoList_empty(t *testing.T) {
	list := models.NewTodoList("Shopping")
	if list.Name != "Shopping" {
		t.Errorf("expected name 'Shopping', got %q", list.Name)
	}
	if list.TotalItems() != 0 {
		t.Errorf("expected 0 items, got %d", list.TotalItems())
	}
}

func TestTodoList_AddAndCheck(t *testing.T) {
	list := models.NewTodoList("Tasks")
	item := list.AddItem("Buy milk")
	if list.TotalItems() != 1 {
		t.Fatalf("expected 1 item, got %d", list.TotalItems())
	}
	if err := list.CheckItem(item.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if list.DoneItems() != 1 {
		t.Errorf("expected 1 done, got %d", list.DoneItems())
	}
}

func TestTodoList_Uncheck(t *testing.T) {
	list := models.NewTodoList("Work")
	item := list.AddItem("Send email")
	_ = list.CheckItem(item.ID)
	if err := list.UncheckItem(item.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if list.DoneItems() != 0 {
		t.Error("expected 0 done items after uncheck")
	}
}

func TestTodoList_Remove(t *testing.T) {
	list := models.NewTodoList("Temp")
	item := list.AddItem("Remove me")
	if err := list.RemoveItem(item.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if list.TotalItems() != 0 {
		t.Error("expected 0 items after removal")
	}
}

func TestTodoList_Progress(t *testing.T) {
	list := models.NewTodoList("Progress")
	a := list.AddItem("A")
	b := list.AddItem("B")
	_ = list.CheckItem(a.ID)
	_ = b
	if list.Progress() != 50.0 {
		t.Errorf("expected 50%%, got %.1f%%", list.Progress())
	}
}

func TestTodoList_SetDescription(t *testing.T) {
	list := models.NewTodoList("My List")
	list.SetDescription("A helpful description")
	if list.Description != "A helpful description" {
		t.Errorf("expected description to be set, got %q", list.Description)
	}
}
