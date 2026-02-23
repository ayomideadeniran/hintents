// Copyright 2025 Erst Users
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/dotandev/hintents/internal/cmd"
)

func main() {
	os.Exit(run(cmd.Execute, os.Stderr))
}

func run(execute func() error, stderr io.Writer) int {
	if err := execute(); err != nil {
		if cmd.IsInterrupted(err) {
			fmt.Fprintln(stderr, "Interrupted. Shutting down...")
			return cmd.InterruptExitCode
		}
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return 1
	}
	return 0
}
