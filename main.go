package main

import (
	"os"

	"lotusaccounts/backend"
	"lotusaccounts/cli"
	"lotusaccounts/httpserver"
)

func main() {
	backend.OpenDb()

	// Run http server if no arguments given
	if len(os.Args) == 1 {
		httpserver.Start()
	} else {
		cli.Cli()
	}

	backend.CloseDb()
}	
