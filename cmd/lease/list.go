package lease

import (
	"context"
	"debugger/pkg/client"
	"debugger/proto"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

var ListLeaseCmd = &cobra.Command{
	Use:     "list",
	Short:   "Listing existing leases",
	Aliases: []string{"ls"},
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
		leases, err := dc.ListLease(ctx, &proto.ListLeaseRequest{})
		if err != nil {
			return err
		}

		// Print output
		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		tbl := table.New("ID", "Namespace", "Deployment", "TTL", "Status")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		for _, l := range leases.Leases {
			tbl.AddRow(l.LeaseId, l.Namespace, l.Deployment, l.Ttl, l.Status)
		}

		tbl.Print()

		return nil
	},
}
