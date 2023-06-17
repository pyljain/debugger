package lease

import (
	"context"
	"debugger/pkg/client"
	"debugger/proto"
	"log"

	"github.com/spf13/cobra"
)

var deploymentName string
var namespace string
var ttl int

var RequestLeaseCmd = &cobra.Command{
	Use:     "request",
	Short:   "Request a lease for a deployment",
	Aliases: []string{"req"},
	RunE: func(cmd *cobra.Command, args []string) error {
		serverAddress, err := cmd.Flags().GetString("addr")
		if err != nil {
			return err
		}

		dc, err := client.New(serverAddress)
		if err != nil {
			return err
		}

		ctx := context.Background()
		res, err := dc.CreateLease(ctx, &proto.CreateLeaseRequest{
			Deployment: deploymentName,
			Namespace:  namespace,
			Ttl:        int32(ttl),
		})
		if err != nil {
			return err
		}

		log.Printf("Requested lease has a Lease ID %d", res.LeaseId)
		return nil
	},
}

func init() {
	RequestLeaseCmd.Flags().StringVar(&deploymentName, "deployment", "", "Please pass the name of the deployment")
	RequestLeaseCmd.Flags().StringVar(&namespace, "namespace", "", "Please pass the namespace")
	RequestLeaseCmd.Flags().IntVar(&ttl, "ttl", 240, "Pass in duration (time to live) in minutes")
}
