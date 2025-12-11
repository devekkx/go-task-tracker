package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.5.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tracker version %s\n", version)
	},
}
