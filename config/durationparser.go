package config

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// We have to implement our own duration parser since go's built in one does
// not handle any units longer than hours, since a day is not always exactly
// 24 hours, (and same with weeks, months, etc). This parser uses an
// approximation of each unit.

type configDuration struct {
	time.Duration
}

var timeUnits = map[string]time.Duration{
	// Exact
	"second": time.Second,
	"minute": time.Minute,
	"hour":   time.Hour,

	// Approximations
	"day":   24 * time.Hour,
	"week":  (24 * time.Hour) * 7,
	"month": (24 * time.Hour) * 30,
	"year":  (24 * time.Hour) * 365,
}

// Removes all empty strings ("") in an array of strings, which can result from double spaces in strings.Split
func skipEmpty(strs []string) []string {
	var result []string
	for _, str := range strs {
		if str != "" {
			result = append(result, str)
		}
	}

	return result
}

type ErrInvalidUnit struct {
	Unit string
}

func (e *ErrInvalidUnit) Error() string {
	return "invalid time unit '" + e.Unit + "'"
}

var (
	ErrDurationMalformed = errors.New("cannot parse malformed duration string")
	ErrDurationOrdering  = errors.New("durations must be given in descending order with no units repeated")
)

// TODO: unit tests
// Format is defined in the default config.toml
func (d *configDuration) UnmarshalText(text []byte) error {
	d.Duration = 0
	// previousUnit is used to ensure that duration strings are given in descending order with no repeated units.
	previousUnit := timeUnits["year"] + time.Second

	str := strings.TrimSpace(string(text))
	sections := skipEmpty(strings.Split(str, " "))

	// We need a minimum of two sections (number and unit) for a duration string
	if len(sections) < 2 {
		return ErrDurationMalformed
	}

	for len(sections) >= 2 {
		num, err := strconv.ParseUint(sections[0], 10, 32)
		if err != nil {
			return err
		}

		// Handle plural units by removing any trailing 's' chars
		unitStr := strings.TrimSuffix(sections[1], "s")

		unitTime, ok := timeUnits[unitStr]
		if !ok {
			return &ErrInvalidUnit{sections[1]}
		}

		// Disallow durations given in non-descending order, and multiple of the same unit
		if unitTime >= previousUnit {
			return ErrDurationOrdering
		}
		previousUnit = unitTime

		d.Duration += time.Duration(num) * unitTime

		sections = sections[2:]
	}

	// The above loop should consume all sections, otherwise there will still be text remaining to be parsed
	if len(sections) != 0 {
		return ErrDurationMalformed
	}

	return nil
}

// Intended to be used only on string literals, as it panics if it is given a value it is unable to parse.
func durationLiteral(str string) configDuration {
	var duration configDuration
	err := duration.UnmarshalText([]byte(str))
	if err != nil {
		panic(err)
	}
	return duration
}
