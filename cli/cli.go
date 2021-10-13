package cli

import (
	"os"
	"log"

	"lotusaccounts/backend"
)

// TODO: replace log in this package with fmt

func cliNewUser() {
	username := cliAskText("username", backend.ValidateNewUsername)
	password := cliAskText("password", backend.ValidateNewPassword)
	is_admin := cliYesNo("is this an admin account?")

	backend.CreateUser(username, password, is_admin)
}

func Cli() {
	switch mode := os.Args[1]; mode {
	case "new-user":
		cliNewUser()	
	case "delete-user":
		log.Fatal("Not implemented")
	case "list-users":
		log.Fatal("Not implemented")
	case "edit-user":
		log.Fatal("Not implemented")
	default:
		log.Fatal("Unknown mode " + mode)
	}
}	
