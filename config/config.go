package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

var Config AccountsConfig

// Documentation for these structs is in the default config.toml
type AccountsConfig struct {
	LoginExpireTime     configDuration
	LoginMaxRefreshTime configDuration
	Domain              string
	// TODO: how to handle relative paths, and is there a service root folder we can have?
	DatabaseLocation string
	Password         PasswordConfig
	Logging          LoggingConfig
}

type PasswordConfig struct {
	MinLength                    int
	MaxLength                    int
	RequireLowerCaseAndUpperCase bool
	RequireSpecialCharacters     bool
	RequireNumber                bool
}

type LoggingConfig struct {
	LogTokenIps   bool
	LogFailedAuth bool

	// NIT: these two fields should really be handled by the logging daemon, as they control log rotation
	MaxFailedAuth       int
	FailedAuthRetention configDuration
}

/*
type Service struct {
	Name string
	Subdomain string
	IconPath string
}
*/

func LoadConfig() error {
	fileData, err := ioutil.ReadFile("config.toml")
	if err != nil {
		return err
	}

	// Default configuration
	Config = AccountsConfig{
		LoginExpireTime:     durationLiteral("1 week"),
		LoginMaxRefreshTime: durationLiteral("6 months"),
		DatabaseLocation:    "auth.db",
		Password: PasswordConfig{
			MinLength:                    8,
			MaxLength:                    256,
			RequireLowerCaseAndUpperCase: true,
			RequireSpecialCharacters:     true,
			RequireNumber:                true,
		},
		Logging: LoggingConfig{
			LogTokenIps:         true,
			LogFailedAuth:       true,
			MaxFailedAuth:       10_000,
			FailedAuthRetention: durationLiteral("1 year"),
		},
	}

	if _, err := toml.Decode(string(fileData), &Config); err != nil {
		// NIT: error formating sucks here, it doesn't give any line info
		return err
	}
	// TODO: check for required fields and constraints on all fields
	if Config.DatabaseLocation == "" {
		return errors.New("missing Domain in config.toml")
	} else if !(Config.Password.MaxLength >= Config.Password.MinLength) {
		return fmt.Errorf(
			"password maximum length (currently %d) be greater than the minimum len (currently %d) defined in config.toml",
			Config.Password.MaxLength,
			Config.Password.MinLength)
	}

	return nil
}
