package backend

import (
	"database/sql"
	"log"

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

var db *sql.DB

func OpenDb() {
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

func CloseDb() {
	db.Close()
}
