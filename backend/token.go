package backend

import (
	"log"
	"errors"
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

var (
	UsernamePasswordMismatch = errors.New("Invalid password for this user")
)

// Returns nil if valid, otherwise returns error
func validateUsernamePassword(username string, password string) error {
	user, err := GetUserBy(QueryTypeUsername, username, true)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) 
	if err == nil {
		return nil
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		return UsernamePasswordMismatch
	} else {
		// Unknown bcrypt error
		log.Panic(err)
	}
	panic("unreachable")
}

func CreateToken(username string, password string, location string) (string, error) {
	err := validateUsernamePassword(username, password)
	if err != nil {
		return "", err
	}

	// generate random base64 string to be used as bearer token
	// a length of 24 bytes is chosen to be "probably good enough"
	tokenBytes := make([]byte, 24)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		log.Panic(err)
	}

	token := base64.StdEncoding.EncodeToString(tokenBytes)

	// TODO: insert into database
	// TODO: read config file for expiry and stuff

	return token, nil
}
