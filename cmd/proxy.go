package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zquestz/ws-tcp-proxy/server"
)

// Stores configuration data.
var config Config

// ProxyCmd is the main command for Cobra.
var ProxyCmd = &cobra.Command{
	Use:   "ws-tcp-proxy <address>",
	Short: "Simple websocket tcp proxy.",
	Long:  `Simple websocket tcp proxy.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := performCommand(cmd, args)
		if err != nil {
			bail(err)
		}
	},
}

func init() {
	err := config.Load()
	if err != nil {
		bail(fmt.Errorf("Failed to load configuration: %s", err))
	}

	prepareFlags()
}

func bail(err error) {
	fmt.Fprintf(os.Stderr, "[Error] %s\n", err)
	os.Exit(1)
}

func prepareFlags() {
	if config.Port == 0 {
		config.Port = defaultPort
	}

	ProxyCmd.PersistentFlags().BoolVarP(
		&config.DisplayVersion, "version", "v", false, "display version")
	ProxyCmd.PersistentFlags().IntVarP(
		&config.Port, "port", "p", config.Port, "server port")
	ProxyCmd.PersistentFlags().StringVarP(
		&config.Cert, "cert", "c", config.Cert, "path to cert.pem for TLS")
	ProxyCmd.PersistentFlags().StringVarP(
		&config.Key, "key", "k", config.Key, "path to key.pem for TLS")
	ProxyCmd.PersistentFlags().BoolVarP(
		&config.TextMode, "text-mode", "t", config.TextMode, "text mode")
}

// Where all the work happens.
func performCommand(cmd *cobra.Command, args []string) error {
	if config.DisplayVersion {
		fmt.Printf("%s %s\n", appName, version)
		return nil
	}

	if len(args) != 1 {
		help := cmd.HelpFunc()
		help(cmd, args)

		return nil
	}

	address := args[0]

	err := server.Run(config.Port, config.Cert, config.Key, config.TextMode, address)
	if err != nil {
		return err
	}

	return nil
}
