package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "1.0.0-rc1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tracker version %s\n", version)
	},
}
