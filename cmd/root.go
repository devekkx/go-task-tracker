package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tracker",
	Short: "A CLI task and todo list manager",
	Long: `tracker is a fast, offline-first CLI tool for managing tasks and todo lists.

Your data is stored at ~/.task-tracker/data.json.

Examples:
  tracker task add "Buy groceries" --priority medium
  tracker task list
  tracker todo create "Weekend Tasks"`,
}

// Execute runs the root command
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable color output")
	rootCmd.AddCommand(taskCmd)
	rootCmd.AddCommand(todoCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(statsCmd)
}
