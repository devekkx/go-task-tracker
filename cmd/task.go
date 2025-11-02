package cmd

import (
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
	Long:  "Add, list, update, and delete tasks.",
}

func init() {
	taskCmd.AddCommand(taskAddCmd)
}
