package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/devekkx/go-task-tracker/internal/display"
	"github.com/devekkx/go-task-tracker/internal/models"
	"github.com/devekkx/go-task-tracker/internal/storage"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
	Long:  "Add, list, update, and delete tasks.",
}

var (
	addDesc     string
	addPriority string
	addDue      string
	addTags     string
)

var taskAddCmd = &cobra.Command{
	Use:   "add <title>",
	Short: "Add a task",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := strings.Join(args, " ")
		priority, err := models.ValidPriority(addPriority)
		if err != nil {
			return err
		}
		task := models.NewTask(title, addDesc, priority)
		if addDue != "" {
			d, err := time.Parse("2006-01-02", addDue)
			if err != nil {
				return fmt.Errorf("invalid due date %q: use YYYY-MM-DD format", addDue)
			}
			task.SetDueDate(d)
		}
		if addTags != "" {
			for _, tag := range strings.Split(addTags, ",") {
				task.AddTag(strings.TrimSpace(tag))
			}
		}
		store, err := storage.New()
		if err != nil {
			return err
		}
		if err := store.AddTask(task); err != nil {
			return err
		}
		display.Success("Task added: %s (ID: %s)", task.Title, task.ID)
		return nil
	},
}

var (
	listStatus   string
	listPriority string
	listTag      string
	listSearch   string
	listArchived bool
	listSort     string
	listLimit    int
	listJSON     bool
)

var taskListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil {
			return err
		}
		opts := storage.FilterOptions{
			Status:   listStatus,
			Priority: listPriority,
			Tag:      listTag,
			Search:   listSearch,
		}
		opts.SortBy = listSort
		if listArchived {
			t := true
			opts.Archived = &t
		}
		tasks := store.ListTasks(opts)
		if listLimit > 0 && len(tasks) > listLimit {
			tasks = tasks[:listLimit]
		}
		if listJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(tasks)
		}
		display.PrintTasks(tasks)
		return nil
	},
}

var taskShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show details of a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil {
			return err
		}
		task, err := store.GetTask(args[0])
		if err != nil {
			return err
		}
		display.PrintTask(task)
		return nil
	},
}

var taskDoneCmd = &cobra.Command{
	Use:   "done <id>",
	Short: "Mark a task as done",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil {
			return err
		}
		task, err := store.GetTask(args[0])
		if err != nil {
			return err
		}
		task.MarkDone()
		if err := store.UpdateTask(task); err != nil {
			return err
		}
		display.Success("Task marked as done: %s", task.Title)
		return nil
	},
}

var taskStartCmd = &cobra.Command{
	Use:   "start <id>",
	Short: "Mark a task as in-progress",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil {
			return err
		}
		task, err := store.GetTask(args[0])
		if err != nil {
			return err
		}
		task.MarkInProgress()
		if err := store.UpdateTask(task); err != nil {
			return err
		}
		display.Success("Task started: %s", task.Title)
		return nil
	},
}

var taskDeleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"rm"},
	Short:   "Delete a task",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil {
			return err
		}
		task, err := store.GetTask(args[0])
		if err != nil {
			return err
		}
		if err := store.DeleteTask(args[0]); err != nil {
			return err
		}
		display.Success("Task deleted: %s", task.Title)
		return nil
	},
}

var (
	updateTitle    string
	updateDesc     string
	updatePriority string
	updateStatus   string
	updateDue      string
	updateTags     string
)

var taskUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil {
			return err
		}
		task, err := store.GetTask(args[0])
		if err != nil {
			return err
		}
		changed := false
		if updateTitle != "" {
			task.Title = updateTitle
			changed = true
		}
		if updateDesc != "" {
			task.Description = updateDesc
			changed = true
		}
		if updatePriority != "" {
			p, err := models.ValidPriority(updatePriority)
			if err != nil {
				return err
			}
			task.Priority = p
			changed = true
		}
		if updateStatus != "" {
			s, err := models.ValidStatus(updateStatus)
			if err != nil {
				return err
			}
			task.Status = s
			changed = true
		}
		if updateDue != "" {
			d, err := time.Parse("2006-01-02", updateDue)
			if err != nil {
				return fmt.Errorf("invalid due date %q: use YYYY-MM-DD", updateDue)
			}
			task.SetDueDate(d)
			changed = true
		}
		if updateTags != "" {
			task.Tags = []string{}
			for _, tag := range strings.Split(updateTags, ",") {
				task.AddTag(strings.TrimSpace(tag))
			}
			changed = true
		}
		if !changed {
			display.Warn("No changes provided.")
			return nil
		}
		if err := store.UpdateTask(task); err != nil {
			return err
		}
		display.Success("Task updated: %s", task.Title)
		return nil
	},
}

var taskCopyCmd = &cobra.Command{
	Use:   "copy <id>",
	Short: "Duplicate a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil {
			return err
		}
		copied, err := store.CopyTask(args[0])
		if err != nil {
			return err
		}
		display.Success("Task copied: %s (ID: %s)", copied.Title, copied.ID)
		return nil
	},
}

var taskBulkDoneCmd = &cobra.Command{
	Use:   "bulk-done",
	Short: "Mark all filtered tasks as done",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil {
			return err
		}
		opts := storage.FilterOptions{
			Status:   listStatus,
			Priority: listPriority,
			Tag:      listTag,
		}
		n, err := store.BulkMarkDone(opts)
		if err != nil {
			return err
		}
		if n == 0 {
			display.Info("No matching tasks.")
		} else {
			display.Success("Marked %d task(s) as done.", n)
		}
		return nil
	},
}

var taskClearDoneCmd = &cobra.Command{
	Use:   "clear-done",
	Short: "Remove all completed tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil {
			return err
		}
		n, err := store.ClearDoneTasks()
		if err != nil {
			return err
		}
		if n == 0 {
			display.Info("No completed tasks to remove.")
		} else {
			display.Success("Removed %d completed task(s).", n)
		}
		return nil
	},
}

var taskArchiveCmd = &cobra.Command{
	Use:   "archive <id>",
	Short: "Archive a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil {
			return err
		}
		if err := store.ArchiveTask(args[0]); err != nil {
			return err
		}
		display.Success("Task archived: %s", args[0])
		return nil
	},
}

var taskUnarchiveCmd = &cobra.Command{
	Use:   "unarchive <id>",
	Short: "Restore an archived task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil {
			return err
		}
		if err := store.UnarchiveTask(args[0]); err != nil {
			return err
		}
		display.Success("Task unarchived: %s", args[0])
		return nil
	},
}

func init() {
	taskAddCmd.Flags().StringVarP(&addDesc, "desc", "d", "", "Task description")
	taskAddCmd.Flags().StringVarP(&addPriority, "priority", "p", "medium", "Priority: low, medium, high")
	taskAddCmd.Flags().StringVar(&addDue, "due", "", "Due date (YYYY-MM-DD)")
	taskAddCmd.Flags().StringVarP(&addTags, "tags", "t", "", "Comma-separated tags")
	taskListCmd.Flags().StringVarP(&listStatus, "status", "s", "", "Filter by status")
	taskListCmd.Flags().StringVarP(&listPriority, "priority", "p", "", "Filter by priority")
	taskListCmd.Flags().StringVarP(&listTag, "tag", "t", "", "Filter by tag")
	taskListCmd.Flags().StringVarP(&listSearch, "search", "q", "", "Search title and description")
	taskListCmd.Flags().BoolVar(&listArchived, "archived", false, "Show archived tasks instead")
	taskListCmd.Flags().StringVarP(&listSort, "sort", "S", "", "Sort by: title, priority, due, created")
	taskListCmd.Flags().IntVarP(&listLimit, "limit", "n", 0, "Limit number of results (0 = no limit)")
	taskListCmd.Flags().BoolVar(&listJSON, "json", false, "Output as JSON")
	taskCmd.AddCommand(taskAddCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskShowCmd)
	taskCmd.AddCommand(taskDoneCmd)
	taskCmd.AddCommand(taskStartCmd)
	taskCmd.AddCommand(taskDeleteCmd)
	taskUpdateCmd.Flags().StringVar(&updateTitle, "title", "", "New title")
	taskUpdateCmd.Flags().StringVarP(&updateDesc, "desc", "d", "", "New description")
	taskUpdateCmd.Flags().StringVarP(&updatePriority, "priority", "p", "", "New priority")
	taskUpdateCmd.Flags().StringVarP(&updateStatus, "status", "s", "", "New status")
	taskUpdateCmd.Flags().StringVar(&updateDue, "due", "", "New due date (YYYY-MM-DD)")
	taskUpdateCmd.Flags().StringVarP(&updateTags, "tags", "t", "", "New tags (replaces existing)")
	taskCmd.AddCommand(taskUpdateCmd)
	taskCmd.AddCommand(taskCopyCmd)
	taskCmd.AddCommand(taskBulkDoneCmd)
	taskCmd.AddCommand(taskClearDoneCmd)
	taskCmd.AddCommand(taskArchiveCmd)
	taskCmd.AddCommand(taskUnarchiveCmd)
}
