package cmd

import (
	"linode-ddns/pkg/client_linode"

	"github.com/spf13/cobra"
)

var (
	record       int
	ip           string
	linodeClient = &cobra.Command{
		Use:   "linode",
		Short: "linode",
		Long:  "linode provides the command line to update domains",
		RunE: func(cmd *cobra.Command, args []string) error {
			return client_linode.Client(cmd.Context(), apiKey, debug, record, ip)
		},
	}
)

func init() {
	rootCmd.AddCommand(linodeClient)

	linodeClient.SuggestionsMinimumDistance = 1

	linodeClient.PersistentFlags().IntVar(&record, "record", 0, "Record ID to update")
	_ = linodeClient.MarkFlagRequired("record")

	linodeClient.PersistentFlags().StringVar(&ip, "ip", "", "IP address for update")
	_ = linodeClient.MarkFlagRequired("ip")
}
