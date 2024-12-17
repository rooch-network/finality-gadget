package main

import (
	"log"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "roochfgd",
		Short:         "roochfgd - Babylon Rooch Finality Gadget",
		Long:          `roochfgd is a daemon to track consecutive quorum and query the Babylon BTC block finalization status of Rooch chains.`,
		SilenceErrors: false,
	}

	return rootCmd
}

func main() {
	cmd := NewRootCmd()

	cmd.AddCommand(CommandStart())

	cmd.PersistentFlags().String("cfg", "config.toml", "config file")
	if err := viper.BindPFlag("cfg", cmd.PersistentFlags().Lookup("cfg")); err != nil {
		log.Fatalf("Error binding flag: %s", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error executing your roochfgd daemon: %s", err)
		os.Exit(1)
	}
}

// Runs cmd with client context and returns an error.
func runEWithClientCtx(
	fRunWithCtx func(ctx client.Context, cmd *cobra.Command, args []string) error,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		clientCtx, err := client.GetClientQueryContext(cmd)
		if err != nil {
			return err
		}

		return fRunWithCtx(clientCtx, cmd, args)
	}
}
