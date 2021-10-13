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
	UsernameNotFound = errors.New("Username not found")
	UsernamePasswordMismatch = errors.New("Invalid password for this user")
	// TODO: add errors for password not meeting constraints
)

// Returns nil if valid, otherwise returns error
func validateUsernamePassword(username string, password string) error {
	// TODO: backend.GetUser(username string)
	rows, err := db.Query(
		"SELECT password_hash FROM users WHERE username = ?", 
		username)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err != nil {
			log.Panic(rows.Err())
		} else {
			return UsernameNotFound
		}
	}

	var dbHash []byte
	err = rows.Scan(&dbHash)
	if err != nil {
		log.Panic(err)
	}

	// TODO: change hashpassword error 
	err = bcrypt.CompareHashAndPassword(dbHash, []byte(password)) 
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

func CreateNewToken(username string, password string, location string) (string, error) {
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
