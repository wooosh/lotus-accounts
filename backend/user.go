package backend

import (
	"fmt"
	"log"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// Wrapper around hashing mechanism
func hashPassword(password string) ([]byte, error) {
	// TODO: prefix with hash scheme
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

var (
	ErrUsernameLength = errors.New("Username must be between 3 and 24 characters long")
	ErrUsernameInvalidChars = errors.New("Username cannot contain non-alphanumeric/underscore characters")
	ErrUsernameAlreadyInUse = errors.New("Username already in use")
	
	usernameAllowedCharsRegex = regexp.MustCompile("^[a-zA-Z0-9_]*$")
)
// TODO: unit tests

func ValidateNewUsername(username string) error {
	// TODO: check if username already exists in DB
	if !(3 <= len(username) && len(username) <= 24) {
		return ErrUsernameLength
	}

	// TODO: only compile once
	if !usernameAllowedCharsRegex.MatchString(username) {
		return ErrUsernameInvalidChars
	}

	_, err := GetUserBy(QueryTypeUsername, username, false)
	if err == nil {
		return ErrUsernameAlreadyInUse
	} else if err != ErrUserNotFound {
		// TODO: send appropriate http code when server panics
		log.Panic(err)
	}

	return nil
}

var (
	ErrPasswordLength = errors.New("Password must be between 8 and 256 characters long")
)

func ValidateNewPassword(password string) error {
	if !(8 <= len(password) && len(password) <= 256) {
		return ErrPasswordLength
	}

	// TODO: determine if more password requirements (require numbers, special
	// chars, capitalization) would be beneficial. maybe have options in
	// the config

	return nil
}

type UserIdType int
const (
	QueryTypeUserId UserIdType = iota
	QueryTypeUsername
)

func uniqueUserIdentifierFieldName(name UserIdType) string {
	switch name {
	case QueryTypeUserId:
		return "id"
	case QueryTypeUsername:
		return "username"
	default:
		panic("Invalid UniqueUserIdType")
	}
}

var ErrUserNotFound = errors.New("No user found with given ID/Username")

func GetUserBy(idType UserIdType, value interface{}, includePasswordHash bool) (User, error) {
	var user User
	
	rows, err := db.Query(
		fmt.Sprintf("SELECT id, username, is_admin, password_hash FROM users WHERE %s = ?", 
			    uniqueUserIdentifierFieldName(idType)), 
		value)
	if err != nil {
		log.Panic(err)
	}

	defer rows.Close()

	if !rows.Next() {
		if err != nil {
			log.Panic(rows.Err())
		} else {
			return user, ErrUserNotFound
		}
	}

	err = rows.Scan(&user.Id, &user.Username, &user.IsAdmin, &user.PasswordHash)
	if err != nil {
		log.Panic(err)
	}

	if !includePasswordHash {
		user.PasswordHash = nil
	}

	return user, nil
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
