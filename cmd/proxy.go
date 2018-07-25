package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zquestz/ws-tcp-proxy/server"
	"golang.org/x/crypto/acme/autocert"
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
	ProxyCmd.PersistentFlags().BoolVarP(
		&config.TCPTLS, "tcp-tls", "", config.TCPTLS, "connect to TCP address via TLS")
	ProxyCmd.PersistentFlags().StringVarP(
		&config.TCPTLSCert, "tcp-tls-cert", "", config.TCPTLSCert, "path to client.crt for TCP TLS")
	ProxyCmd.PersistentFlags().StringVarP(
		&config.TCPTLSKey, "tcp-tls-key", "", config.TCPTLSKey, "path to client.key for TCP TLS")
	ProxyCmd.PersistentFlags().StringVarP(
		&config.TCPTLSRootCA, "tcp-tls-root-ca", "", config.TCPTLSRootCA, "path to ca.crt for TCP TLS")
	ProxyCmd.PersistentFlags().StringVarP(
		&config.AutoCert, "auto-cert", "a", config.AutoCert, "register hostname with LetsEncrypt")
}

// Where all the work happens.
func performCommand(cmd *cobra.Command, args []string) error {
	if config.DisplayVersion {
		fmt.Printf("%s %s\n", appName, version)
		return nil
	}

	if config.AutoCert != "" && (config.Cert != "" || config.Key != "") {
		return errors.New("can't specify auto-cert and key/cert")
	}

	if len(args) != 1 {
		help := cmd.HelpFunc()
		help(cmd, args)

		return nil
	}

	address := args[0]

	m, err := getLetsEncryptManager()
	if err != nil {
		return err
	}

	serverConfig := &server.Config{
		Port:         config.Port,
		Cert:         config.Cert,
		Key:          config.Key,
		TextMode:     config.TextMode,
		Address:      address,
		TCPTLS:       config.TCPTLS,
		TCPTLSRootCA: config.TCPTLSRootCA,
		TCPTLSCert:   config.TCPTLSCert,
		TCPTLSKey:    config.TCPTLSKey,
		CertManager:  m,
	}

	err = server.Run(serverConfig)
	if err != nil {
		return err
	}

	return nil
}

func getLetsEncryptManager() (*autocert.Manager, error) {
	configDir, err := config.configDir()
	if err != nil {
		return nil, err
	}

	certDir := filepath.Join(configDir, "certs")
	if _, err := os.Stat(certDir); os.IsNotExist(err) {
		err = os.MkdirAll(certDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	if config.AutoCert == "" {
		return nil, nil
	}

	m := &autocert.Manager{
		Cache:      autocert.DirCache(certDir),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.AutoCert),
	}

	go http.ListenAndServe(":http", m.HTTPHandler(nil))

	return m, nil
}
