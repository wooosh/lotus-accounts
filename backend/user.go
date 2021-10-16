package backend

import (
	"errors"
	"fmt"
	"log"
	"lotusaccounts/config"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Wrapper around hashing mechanism
func hashPassword(password string) ([]byte, error) {
	// TODO: prefix with hash scheme
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

var (
	ErrUsernameLength       = errors.New("username must be between 3 and 24 characters long")
	ErrUsernameInvalidChars = errors.New("username cannot contain non-alphanumeric/underscore characters")
	ErrUsernameAlreadyInUse = errors.New("username already in use")

	usernameAllowedCharsRegex = regexp.MustCompile("^[a-zA-Z0-9_]*$")
)

// TODO: unit tests

// Checks if a username is a valid name for a new account
// TODO: TOCTOU error
func ValidateNewUsername(username string) error {
	// Ensure username is a reasonable length
	if !(3 <= len(username) && len(username) <= 24) {
		return ErrUsernameLength
	}

	// Ensure the username is composed of only alphanumeric chars and underscores
	if !usernameAllowedCharsRegex.MatchString(username) {
		return ErrUsernameInvalidChars
	}

	// Ensure username is not already taken
	_, err := GetUserBy(QueryTypeUsername, username, false)
	if err == nil {
		return ErrUsernameAlreadyInUse
	} else if err != ErrUserNotFound {
		// TODO: send appropriate http code and message when server panics
		log.Panic(err)
	}

	return nil
}

var (
	ErrPasswordLength              = errors.New("password must be between 8 and 256 characters long")
	ErrPasswordLacksSpecialChars   = errors.New("password must have atleast one of the following special characters:  !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~")
	ErrPasswordLacksCapitalization = errors.New("password must have atleast one lowercase and one uppercase letter")
	ErrPasswordLacksNumber         = errors.New("password must have atleast one number")

	uppercaseRegex = regexp.MustCompile("[[:upper:]]")
	lowercaseRegex = regexp.MustCompile("[[:lower:]]")
)

// TODO: unit test
// Check if a password meets the requirements defined in the config
func ValidateNewPassword(password string) error {
	// Ensure password is between length limits
	if !(config.Config.Password.MinLength <= len(password) &&
		len(password) <= config.Config.Password.MaxLength) {
		return ErrPasswordLength
	}

	// Ensure password contains special characters if it is required by the config
	if config.Config.Password.RequireSpecialCharacters &&
		!strings.ContainsAny(password, " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~") {
		return ErrPasswordLacksSpecialChars
	}

	// Ensure password contains atleast one number if required by the config
	if config.Config.Password.RequireLowerCaseAndUpperCase &&
		!strings.ContainsAny(password, "0123456789") {
		return ErrPasswordLacksNumber
	}

	// Ensure password contains one uppercase and one lowercase leter if required by the config
	if config.Config.Password.RequireLowerCaseAndUpperCase &&
		!(uppercaseRegex.MatchString(password) && lowercaseRegex.MatchString(password)) {
		return ErrPasswordLacksCapitalization
	}

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

var ErrUserNotFound = errors.New("no user found with given ID/Username")

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
