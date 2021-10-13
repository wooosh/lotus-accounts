package main

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Wrapper around hashing mechanism
func hashPassword(password string) ([]byte, error) {
	// TODO: prefix with hash scheme
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func validateNewUsername(username string) error {
	// TODO: check if username already exists in DB
	if len(username) <= 2 {
		return errors.New("Username must be greater than 2 characters long")
	}

	return nil
}

func validateNewPassword(password string) error {
	if len(password) < 8 {
		return errors.New("Password must be atleast 8 characters long")
	}

	return nil
}
