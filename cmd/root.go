package cmd

import (
	"fmt"
	"linode-ddns/pkg/client_linode"
	"os"

	"github.com/spf13/cobra"
)

var (
	apiKey      string
	domains     bool
	debug       bool
	refreshTime int

	rootCmd = &cobra.Command{
		Use:   "linode-ddns",
		Short: "linode-ddns provides a client_linode to update domains",
		Long: `linode-ddns provides a client_linode to update domains and 
also can run in daemon mode to update the record automatically`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if domains {
				return client_linode.Client(cmd.Context(), apiKey, debug, 0, "")
			}
			return nil
		},
	}
	daemonClient = &cobra.Command{
		Use:   "daemon",
		Short: "daemon",
		Long:  "Runs in daemon mode to keep IP address updated",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.SuggestionsMinimumDistance = 1
	rootCmd.PersistentFlags().StringVar(&apiKey, "apiKey", "", "linode apiKey (required)")
	_ = rootCmd.MarkFlagRequired("apiKey")
	rootCmd.PersistentFlags().BoolVarP(&domains, "domains", "d", false, "List the domains in the linode account (default)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug on the Linode client")

	daemonClient.SuggestionsMinimumDistance = 1
	daemonClient.PersistentFlags().IntVarP(&refreshTime, "refresh-time", "r", 5, "Refresh time for the daemon (minutes)") //nolint
	rootCmd.AddCommand(daemonClient)
}

func initConfig() {
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Needs to specify a apiKey")
		os.Exit(1)
	}
}
