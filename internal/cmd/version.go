package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.1.0-alpha"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of erst",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("erst version %s\n", version)
	},
}
