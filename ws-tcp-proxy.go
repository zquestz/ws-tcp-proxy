package main

import (
	"fmt"
	"os"

	"github.com/zquestz/ws-tcp-proxy/cmd"
)

func main() {
	setupSignalHandlers()

	if err := cmd.ProxyCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
