package cli

import (
	"fmt"
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
