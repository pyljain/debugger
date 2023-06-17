package cmd

import (
	"debugger/cmd/lease"

	"github.com/spf13/cobra"
)

var leaseCmd = &cobra.Command{
	Use:   "lease",
	Short: "Managing lease life cycle",
}

func init() {
	leaseCmd.AddCommand(lease.ListLeaseCmd)
	leaseCmd.AddCommand(lease.RequestLeaseCmd)
	leaseCmd.AddCommand(lease.ApproveLeaseCmd)
	leaseCmd.PersistentFlags().String("addr", "localhost:9091", "the address to connect to")
}
