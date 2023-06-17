package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "debugger",
	Short: "Debugger is a very fast CLI to debug running pods in Kubernetes",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(leaseCmd)
	rootCmd.AddCommand(serverCmd)
}
