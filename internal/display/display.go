package display

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"

	"github.com/devekkx/go-task-tracker/internal/models"
	"github.com/devekkx/go-task-tracker/internal/storage"
)

var (
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	warnColor    = color.New(color.FgYellow, color.Bold)
	headerColor  = color.New(color.FgCyan, color.Bold)
	dimColor     = color.New(color.Faint)
	boldColor    = color.New(color.Bold)
)

// Success renders a success message
// Success prints a success message
func Success(format string, args ...any) {
	successColor.Fprintf(os.Stdout, "✓ "+format+"\n", args...)
}

// Error prints an error message to stderr.
func Error(format string, args ...any) {
	errorColor.Fprintf(os.Stderr, "✗ "+format+"\n", args...)
}

// Warn prints a warning message.
func Warn(format string, args ...any) {
	warnColor.Fprintf(os.Stdout, "⚠ "+format+"\n", args...)
}

// Info prints an informational message.
func Info(format string, args ...any) {
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}

// PrintTasks renders a table of tasks.
func PrintTasks(tasks []models.Task) {
	if len(tasks) == 0 {
		dimColor.Println("No tasks found.")
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Priority", "Status", "Due Date", "Tags"})
	table.SetBorder(false)
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
	)
	for _, t := range tasks {
		row := []string{
			truncate(t.ID, 16),
			truncate(t.Title, 40),
			priorityStr(t.Priority),
			statusStr(t.Status),
			dueDateStr(t.DueDate, t.IsOverdue()),
			strings.Join(t.Tags, ", "),
		}
		table.Rich(row, rowColors(t))
	}
	table.Render()
	fmt.Printf("\nTotal: %d task(s)\n", len(tasks))
}

// PrintTask renders a detailed view of a single task.
func PrintTask(t *models.Task) {
	headerColor.Printf("\n  Task: %s\n", t.Title)
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("  %-14s %s\n", "ID:", t.ID)
	fmt.Printf("  %-14s %s\n", "Status:", statusStr(t.Status))
	fmt.Printf("  %-14s %s\n", "Priority:", priorityStr(t.Priority))
	if t.Description != "" {
		fmt.Printf("  %-14s %s\n", "Description:", t.Description)
	}
	if t.DueDate != nil {
		due := t.DueDate.Format("2006-01-02")
		if t.IsOverdue() {
			due = warnColor.Sprint(due + " (overdue)")
		}
		fmt.Printf("  %-14s %s\n", "Due Date:", due)
	}
	if len(t.Tags) > 0 {
		fmt.Printf("  %-14s %s\n", "Tags:", strings.Join(t.Tags, ", "))
	}
	fmt.Printf("  %-14s %s\n", "Created:", t.CreatedAt.Format("2006-01-02 15:04"))
	fmt.Printf("  %-14s %s\n", "Updated:", t.UpdatedAt.Format("2006-01-02 15:04"))
	fmt.Println()
}

// PrintTodoLists renders a table of todo lists.
func PrintTodoLists(lists []models.TodoList) {
	if len(lists) == 0 {
		dimColor.Println("No todo lists found.")
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Items", "Done", "Progress"})
	table.SetBorder(false)
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
	)
	for _, l := range lists {
		table.Append([]string{
			truncate(l.ID, 16),
			truncate(l.Name, 30),
			fmt.Sprintf("%d", l.TotalItems()),
			fmt.Sprintf("%d", l.DoneItems()),
			fmt.Sprintf("%.0f%%", l.Progress()),
		})
	}
	table.Render()
	fmt.Printf("\nTotal: %d list(s)\n", len(lists))
}

// PrintTodoList renders a single todo list with its items.
func PrintTodoList(l *models.TodoList) {
	headerColor.Printf("\n  List: %s\n", l.Name)
	fmt.Printf("  ID: %s\n", l.ID)
	fmt.Printf("  Progress: %.0f%% (%d/%d done)\n", l.Progress(), l.DoneItems(), l.TotalItems())
	fmt.Println(strings.Repeat("─", 50))
	if len(l.Items) == 0 {
		dimColor.Println("  No items yet.")
	} else {
		for _, item := range l.Items {
			var box string
			if item.Done {
				box = successColor.Sprint("[✓]")
			} else {
				box = "[ ]"
			}
			content := item.Content
			if item.Done {
				content = dimColor.Sprint(content)
			}
			fmt.Printf("  %s %-18s %s\n", box, truncate(item.ID, 18), content)
		}
	}
	fmt.Println()
}

// PrintStats renders aggregate statistics.
func PrintStats(stats storage.Stats) {
	headerColor.Println("\n  Task Tracker Statistics")
	fmt.Println(strings.Repeat("─", 40))
	boldColor.Println("\n  Tasks:")
	fmt.Printf("    %-20s %d\n", "Total:", stats.TotalTasks)
	fmt.Printf("    %-20s %d\n", "Pending:", stats.PendingTasks)
	fmt.Printf("    %-20s %d\n", "In Progress:", stats.InProgressTasks)
	fmt.Printf("    %-20s %d\n", "Done:", stats.DoneTasks)
	if stats.OverdueTasks > 0 {
		fmt.Printf("    %-20s ", "Overdue:")
		warnColor.Printf("%d\n", stats.OverdueTasks)
	}
	boldColor.Println("\n  Todo Lists:")
	fmt.Printf("    %-20s %d\n", "Total Lists:", stats.TotalTodoLists)
	fmt.Printf("    %-20s %d\n", "Total Items:", stats.TotalTodoItems)
	fmt.Printf("    %-20s %d\n", "Done Items:", stats.DoneTodoItems)
	fmt.Printf("    %-20s %d\n", "Pending Items:", stats.TotalTodoItems-stats.DoneTodoItems)
	fmt.Println()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func priorityStr(p models.Priority) string {
	switch p {
	case models.PriorityHigh:
		return color.RedString("high")
	case models.PriorityMedium:
		return color.YellowString("medium")
	case models.PriorityLow:
		return color.GreenString("low")
	default:
		return string(p)
	}
}

func statusStr(s models.Status) string {
	switch s {
	case models.StatusDone:
		return color.GreenString("done")
	case models.StatusInProgress:
		return color.YellowString("in-progress")
	case models.StatusPending:
		return color.CyanString("pending")
	default:
		return string(s)
	}
}

func dueDateStr(d *time.Time, overdue bool) string {
	if d == nil {
		return "-"
	}
	s := d.Format("2006-01-02")
	if overdue {
		return color.YellowString(s)
	}
	return s
}

func rowColors(t models.Task) []tablewriter.Colors {
	var c tablewriter.Colors
	switch t.Status {
	case models.StatusDone:
		c = tablewriter.Colors{tablewriter.FgGreenColor}
	case models.StatusInProgress:
		c = tablewriter.Colors{tablewriter.FgYellowColor}
	default:
		c = tablewriter.Colors{}
	}
	return []tablewriter.Colors{c, c, c, c, c, c}
}
