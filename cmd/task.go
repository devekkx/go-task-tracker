package cmd

import (
	"fmt"
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
	Short: "Add a new task",
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
)

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil {
			return err
		}
		tasks := store.ListTasks(storage.FilterOptions{
			Status:   listStatus,
			Priority: listPriority,
			Tag:      listTag,
			Search:   listSearch,
		})
		display.PrintTasks(tasks)
		return nil
	},
}

var taskShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show task details",
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

func init() {
	taskAddCmd.Flags().StringVarP(&addDesc, "desc", "d", "", "Task description")
	taskAddCmd.Flags().StringVarP(&addPriority, "priority", "p", "medium", "Priority: low, medium, high")
	taskAddCmd.Flags().StringVar(&addDue, "due", "", "Due date (YYYY-MM-DD)")
	taskAddCmd.Flags().StringVarP(&addTags, "tags", "t", "", "Comma-separated tags")
	taskListCmd.Flags().StringVarP(&listStatus, "status", "s", "", "Filter by status")
	taskListCmd.Flags().StringVarP(&listPriority, "priority", "p", "", "Filter by priority")
	taskListCmd.Flags().StringVarP(&listTag, "tag", "t", "", "Filter by tag")
	taskListCmd.Flags().StringVarP(&listSearch, "search", "q", "", "Search title and description")
	taskCmd.AddCommand(taskAddCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskShowCmd)
}
