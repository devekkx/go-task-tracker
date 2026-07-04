package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is the current application release.
const Version = "1.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tracker version %s\n", Version)
	},
}
