package lease

import (
	"context"
	"debugger/pkg/client"
	"debugger/proto"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

var ApproveLeaseCmd = &cobra.Command{
	Use:     "approve",
	Short:   "Approve a lease for a deployment",
	Aliases: []string{"apr"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		leaseId := args[0]
		// Call the gRPC client function for ApproveLease
		serverAddress, err := cmd.Flags().GetString("addr")
		if err != nil {
			return err
		}
		c, err := client.New(serverAddress)
		if err != nil {
			return err
		}

		ctx := context.Background()
		leaseIdInt, err := strconv.Atoi(leaseId)
		if err != nil {
			return err
		}

		lease, err := c.ApproveLease(ctx, &proto.ApproveLeaseRequest{
			LeaseId: int32(leaseIdInt),
		})
		if err != nil {
			return err
		}

		// Query the Leases table with a Lease ID
		log.Printf("The lease is now approved %d", lease.LeaseId)
		return nil
	},
}
