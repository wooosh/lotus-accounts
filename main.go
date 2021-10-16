package main

import (
	"os"

	"lotusaccounts/backend"
	"lotusaccounts/cli"
	"lotusaccounts/config"
	"lotusaccounts/httpserver"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}

	backend.OpenDb()

	// TODO: show help if no arguments
	// Run http server if no arguments given
	if len(os.Args) == 1 {
		httpserver.Start()
	} else {
		cli.Cli()
	}

	backend.CloseDb()
}
