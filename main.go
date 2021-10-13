package main

import (
	"os"
)

func main() {
	openDb()

	// Handle cli args
	if len(os.Args) == 1 {
		httpServer()
	} else {
		handleCli()
	}

	closeDb()
}	
