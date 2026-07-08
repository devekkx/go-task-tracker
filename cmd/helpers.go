package cmd

import (
	"github.com/spf13/cobra"

	"github.com/devekkx/go-task-tracker/internal/storage"
)

// withStore wraps a command handler that needs an opened Store, removing the
// repeated storage.New()/error-check boilerplate from every RunE.
func withStore(fn func(store *storage.Store, args []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil {
			return err
		}
		return fn(store, args)
	}
}
