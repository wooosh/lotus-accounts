package backend

import (
	"log"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Wrapper around hashing mechanism
func hashPassword(password string) ([]byte, error) {
	// TODO: prefix with hash scheme
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func ValidateNewUsername(username string) error {
	// TODO: check if username already exists in DB
	if len(username) <= 2 {
		return errors.New("Username must be greater than 2 characters long")
	}

	return nil
}

func ValidateNewPassword(password string) error {
	if len(password) < 8 {
		return errors.New("Password must be atleast 8 characters long")
	}

	return nil
}

func CreateUser(username string, password string, is_admin bool) error {
	err := ValidateNewUsername(username)
	if err != nil {
		return err
	}
	
	err = ValidateNewPassword(password)
	if err != nil {
		return err
	}


	hash, err := hashPassword(password)
	if err != nil {
		log.Panic(err)
	}

	_, err = db.Exec(
		`INSERT INTO users (username, password_hash, is_admin)
		VALUES (?, ?, ?)`,
		username, hash, is_admin)
	if err != nil {
		log.Panic(err)
	}

	return nil
}
