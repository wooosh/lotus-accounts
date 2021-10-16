package backend

import (
	"database/sql"
	"log"
	"lotusaccounts/config"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id       uint64
	Username string
	IsAdmin  bool
	// Will be set to "" unless specifically requested
	PasswordHash []byte
}

type Token struct {
	Id     uint64
	UserId uint64

	// Base64 string
	Token string

	// String representing the device a token was issued for, typically
	// composed of IP and user agent string, or a user provided string
	// like "Linux Desktop" or "Samsung Galaxy S4"
	Device string

	// Domain the token is issued for
	Domain string

	Created time.Time
	Expiry  time.Time
}

const (
	// TODO: set up unique properly
	// TODO: CHECK() constraints where applicable
	// TODO: set expiry -1 when force revoked
	// TODO: make it so that getting tokens does not expose the token unless
	// explicitly asked for similar to how hashes work
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
		device		TEXT NOT NULL,
		domain		TEXT NOT NULL,
		created		INTEGER NOT NULL,
		expiry 		INTEGER NOT NULL,
		user 		INTEGER NOT NULL,

		FOREIGN KEY(user) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS token_uses (
		id		INTEGER PRIMARY KEY AUTOINCREMENT,
		token 		INTEGER NOT NULL,
		time		INTEGER NOT NULL,
		ip		TEXT,
		user_agent	TEXT,
		success 	INTEGER NOT NULL,

		FOREIGN KEY(token) REFERENCES tokens(id)
	);

	CREATE TABLE IF NOT EXISTS failed_auth (
		id		INTEGER PRIMARY KEY AUTOINCREMENT,
		time 		INTEGER NOT NULL,
		username	TEXT,
		ip		TEXT,
		user_agent	TEXT
	);
	`
)

var db *sql.DB

func OpenDb() {
	var err error
	db, err = sql.Open("sqlite3", config.Config.DatabaseLocation)
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
