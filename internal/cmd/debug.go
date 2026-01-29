package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dotandev/hintents/internal/rpc"
	"github.com/dotandev/hintents/internal/simulator"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	network  string
	verbose  bool
	wasmPath string
	args     []string
)

var debugCmd = &cobra.Command{
	Use:   "debug <transaction-hash>",
	Short: "Debug a failed Stellar transaction",
	Long: `Fetch and analyze a failed Stellar smart contract transaction.

This command retrieves the transaction envelope from the Stellar network
and provides detailed information about why the transaction failed.

Example:
  erst debug abc123def456... --network testnet
  erst debug abc123def456... --network mainnet --verbose
  erst debug --wasm ./contract.wasm --args '["arg1", "arg2"]'`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDebug,
}

func init() {
	debugCmd.Flags().StringVarP(&network, "network", "n", "testnet", "Stellar network to use (testnet, mainnet)")
	debugCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	debugCmd.Flags().StringVar(&wasmPath, "wasm", "", "Path to local WASM file for local replay (no network required)")
	debugCmd.Flags().StringSliceVar(&args, "args", []string{}, "Mock arguments for local replay (JSON array of strings)")
}

func runDebug(cmd *cobra.Command, cmdArgs []string) error {
	// Local WASM replay mode
	if wasmPath != "" {
		return runLocalWasmReplay()
	}

	// Network transaction replay mode
	if len(cmdArgs) == 0 {
		return fmt.Errorf("transaction hash is required when not using --wasm flag")
	}

	txHash := cmdArgs[0]
	return runNetworkReplay(txHash)
}

func runLocalWasmReplay() error {
	color.Yellow("⚠️  WARNING: Using Mock State (not mainnet data)")
	fmt.Println()

	// Verify WASM file exists
	if _, err := os.Stat(wasmPath); os.IsNotExist(err) {
		return fmt.Errorf("WASM file not found: %s", wasmPath)
	}

	color.Cyan("🔧 Local WASM Replay Mode")
	fmt.Printf("WASM File: %s\n", wasmPath)
	fmt.Printf("Arguments: %v\n", args)
	fmt.Println()

	// Create simulator runner
	runner, err := simulator.NewRunner()
	if err != nil {
		return fmt.Errorf("failed to initialize simulator: %w", err)
	}

	// Create simulation request with local WASM
	req := &simulator.SimulationRequest{
		EnvelopeXdr:   "",  // Empty for local replay
		ResultMetaXdr: "",  // Empty for local replay
		LedgerEntries: nil, // Mock state will be generated
		WasmPath:      &wasmPath,
		MockArgs:      &args,
	}

	// Run simulation
	color.Green("▶ Executing contract locally...")
	resp, err := runner.Run(req)
	if err != nil {
		color.Red("✗ Execution failed: %v", err)
		return err
	}

	// Display results
	fmt.Println()
	color.Green("✓ Execution completed successfully")
	fmt.Println()

	if len(resp.Logs) > 0 {
		color.Cyan("📋 Logs:")
		for _, log := range resp.Logs {
			fmt.Printf("  %s\n", log)
		}
		fmt.Println()
	}

	if len(resp.Events) > 0 {
		color.Cyan("📡 Events:")
		for _, event := range resp.Events {
			fmt.Printf("  %s\n", event)
		}
		fmt.Println()
	}

	if verbose {
		color.Cyan("🔍 Full Response:")
		jsonBytes, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Println(string(jsonBytes))
	}

	return nil
}

func runNetworkReplay(txHash string) error {
	color.Cyan("🌐 Network Transaction Replay Mode")
	fmt.Printf("Transaction Hash: %s\n", txHash)
	fmt.Printf("Network: %s\n", network)
	fmt.Println()

	// Initialize RPC client
	client, err := rpc.NewClient(network)
	if err != nil {
		return fmt.Errorf("failed to initialize RPC client: %w", err)
	}

	// Fetch transaction
	color.Green("▶ Fetching transaction from network...")
	envelope, resultMeta, err := client.GetTransaction(txHash)
	if err != nil {
		color.Red("✗ Failed to fetch transaction: %v", err)
		return err
	}

	if verbose {
		fmt.Printf("Envelope XDR length: %d bytes\n", len(envelope))
		fmt.Printf("ResultMeta XDR length: %d bytes\n", len(resultMeta))
		fmt.Println()
	}

	// Create simulator runner
	runner, err := simulator.NewRunner()
	if err != nil {
		return fmt.Errorf("failed to initialize simulator: %w", err)
	}

	// Create simulation request
	req := &simulator.SimulationRequest{
		EnvelopeXdr:   envelope,
		ResultMetaXdr: resultMeta,
		LedgerEntries: nil, // TODO: Fetch ledger entries
	}

	// Run simulation
	color.Green("▶ Replaying transaction locally...")
	resp, err := runner.Run(req)
	if err != nil {
		color.Red("✗ Replay failed: %v", err)
		return err
	}

	// Display results
	fmt.Println()
	color.Green("✓ Replay completed successfully")
	fmt.Println()

	if len(resp.Logs) > 0 {
		color.Cyan("📋 Logs:")
		for _, log := range resp.Logs {
			fmt.Printf("  %s\n", log)
		}
		fmt.Println()
	}

	if len(resp.Events) > 0 {
		color.Cyan("📡 Events:")
		for _, event := range resp.Events {
			fmt.Printf("  %s\n", event)
		}
		fmt.Println()
	}

	if verbose {
		color.Cyan("🔍 Full Response:")
		jsonBytes, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Println(string(jsonBytes))
	}

	return nil
}
