package config

import (
	"testing"
	"time"
)

const (
	second = time.Second
	minute = time.Minute
	hour   = time.Hour
	day    = time.Hour * 24
	week   = day * 7
	month  = day * 30
	year   = day * 365
)

func assertDurationEquals(t *testing.T, durationStr string, expected time.Duration) {
	var parsedDuration configDuration
	parsedDuration.UnmarshalText([]byte(durationStr))

	if expected != parsedDuration.Duration {
		t.Fatalf(`expected time value of %s, recieved %s instead`, expected, parsedDuration)
	}
}

func assertDurationErr(t *testing.T, durationStr string) {
	var parsedDuration configDuration
	err := parsedDuration.UnmarshalText([]byte(durationStr))

	if err == nil {
		t.Fatalf(`expected an error, recieved %s instead`, err)
	}
}

func TestTimeUnits(t *testing.T) {
	assertDurationEquals(t,
		"1 year 2 months 3 weeks 4 days 5 hours 6 minutes 7 seconds",
		year+2*month+3*week+4*day+5*hour+6*minute+7*second)

	assertDurationEquals(t,
		"3 weeks 4 days 5 hours 6 minutes 7 seconds",
		3*week+4*day+5*hour+6*minute+7*second)

	assertDurationEquals(t,
		"5 hours 6 minutes 7 seconds",
		5*hour+6*minute+7*second)

	assertDurationEquals(t, "7 seconds", 7*second)
}

func TestErrors(t *testing.T) {
	// Incorrect order
	assertDurationErr(t, "4 seconds 2 days")

	// Missing space
	assertDurationErr(t, "2months")

	// Repeated units
	assertDurationErr(t, "1 year 1 year 2 months")

	// Nothing
	assertDurationErr(t, "")

	// Missing units
	assertDurationErr(t, "25")

	// Invalid number
	assertDurationErr(t, "2k seconds")
}
