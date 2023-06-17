package cmd

import (
	"debugger/pkg/server"
	"log"

	"github.com/spf13/cobra"
)

var (
	connectionString string
	kubeconfig       string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Bring the debugger service up",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Debugger service is up")
		err := server.Start(connectionString, kubeconfig)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	serverCmd.Flags().StringVar(&connectionString, "connection-string", "", "Pass the database connection string")
	serverCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Location of the kubeconfig")
}
