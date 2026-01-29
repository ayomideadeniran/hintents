package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "erst",
	Short: "Erst is a specialized developer tool for the Stellar network",
	Long: `Erst is a specialized developer tool for the Stellar network,
designed to solve the "black box" debugging experience on Soroban.

It helps clarify why a Stellar smart contract transaction failed by:
  - Fetching failed transaction envelopes and ledger state
  - Re-executing transactions locally for detailed analysis
  - Mapping execution failures back to readable source code`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Register commands
	rootCmd.AddCommand(debugCmd)
	rootCmd.AddCommand(versionCmd)
}
