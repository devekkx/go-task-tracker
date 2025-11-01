package cmd

import (
	"github.com/spf13/cobra"

	"github.com/devekkx/go-task-tracker/internal/display"
	"github.com/devekkx/go-task-tracker/internal/storage"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show task and todo statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil {
			return err
		}
		stats := store.GetStats()
		display.PrintStats(stats)
		return nil
	},
}
