package main

import (
	"fmt"
	"os"
	"log"
	"strings"
)


func cliYesNo(prompt string) bool {
	for {
		var input string
		fmt.Print(prompt + "[y/n]: ")
		fmt.Scanln(&input)
		switch strings.ToLower(input) {
		case "y", "yes", "true":
			return true
		case "n", "no", "false":
			return false
		default:
			fmt.Println("incorrect format, expected one of 'y', 'yes', 'true', or 'n', 'no', 'false'")
		}
	}
}

// the user will be asked for input until the verifier returns nil
func cliAskText(prompt string, verifier func(input string) error) string {
	for {
		var input string
		fmt.Print(prompt + ": ")
		fmt.Scanln(&input)
		message := verifier(input)
		if message != nil {
			fmt.Println(message)
		} else {
			return input
		}
	}
}

func cliNewUser() {
	username := cliAskText("username", validateNewUsername)
	password := cliAskText("password", validateNewPassword)
	is_admin := cliYesNo("is this an admin account?")

	dbNewUser(username, password, is_admin)
}

func handleCli() {
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
