package main

import (
	"database/sql"
	"log"
	"errors"
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// TODO: set up unique properly
	createTablesStmt = `
	CREATE TABLE IF NOT EXISTS users (
		id 		INTEGER PRIMARY KEY AUTOINCREMENT,
		username 	TEXT NOT NULL UNIQUE,
		password_hash	BLOB NOT NULL,
		is_admin	INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS tokens (
		id		INTEGER PRIMARY KEY AUTOINCREMENT,
		token		TEXT UNIQUE NOT NULL,
		location	TEXT NOT NULL,
		subdomain	TEXT NOT NULL,
		created		INTEGER NOT NULL,
		expiry 		INTEGER NOT NULL,
		user 		INTEGER NOT NULL,
		FOREIGN KEY(user) REFERENCES users(id)
	)

	`

	// TODO: default expiry time
)

var (
	UsernameNotFound = errors.New("Username not found")
	UsernamePasswordMismatch = errors.New("Invalid password for this user")
	// TODO: add errors for password not meeting constraints
)

var db *sql.DB

func dbNewUser(username string, password string, is_admin bool) error {
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`INSERT INTO users (username, password_hash, is_admin)
		VALUES (?, ?, ?)`,
		username, hash, is_admin)
	// TODO: change sql error to log.Panic or something
	if err != nil {
		return err
	}

	return nil
}

// Returns nil if valid, otherwise returns error
func validateUsernamePassword(username string, password string) error {
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

// TODO: rename this to be its own backend package instead of db
// TODO: do above and move validate stuff in here
func createNewToken(username string, password string, location string) (string, error) {
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

func openDb() {
	var err error
	db, err = sql.Open("sqlite3", "auth.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createTablesStmt)
	if err != nil {
		log.Fatal(err)
	}

}

func closeDb() {
	db.Close()
}
