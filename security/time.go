package security

import (
	"encoding/json"
	"fmt"
	"time"
)

const rfc3339Extended = "2006-01-02T15:04:05.999999999-0700"

var formats = []string{
	"2006-01-02 15:04:05.999999999 -0700 MST",
	"2006-01-2 15:04:05",
	"2006-01-2 15:04",
	"2006-01-2 15:04:05 -0700",
	"2006-01-2 15:04 -0700",
	"2006-01-2 15:04:05 -07:00",
	"2006-01-2 15:04 -07:00",
	"2006-01-2 15:04:05 MST",
	"2006-01-2 15:04 MST",
	time.RFC3339,
	rfc3339Extended,
	time.RFC3339Nano,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.UnixDate,
	time.RubyDate,
	time.ANSIC,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
}

// Time represents a Composer-like date
type Time time.Time

// Format proxifies call to time.Time.Format
func (t Time) Format(layout string) string {
	return time.Time(t).Format(layout)
}

// UnmarshalYAML parses a Composer-like date from YAML to a Go time.Time
func (t *Time) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var date string
	var err error
	if err := unmarshal(&date); err != nil {
		return err
	}
	*t, err = parseDate(date)
	return err
}

// UnmarshalJSON parses a Composer-like date from JSON to a Go time.Time
func (t *Time) UnmarshalJSON(data []byte) error {
	var err error
	var date string
	if err := json.Unmarshal(data, &date); err != nil {
		return err
	}
	*t, err = parseDate(date)
	return err
}

// MarshalJSON dumps a Composer-like date to JSON from a Go time.Time
func (t Time) MarshalJSON() ([]byte, error) {
	return time.Time(t).MarshalJSON()
}

func parseDate(date string) (Time, error) {
	if date == "" {
		// far away in the future, means no fix available
		date = "2123-01-02"
	}

	if tt, ok := TryParseTime(date); ok {
		return Time(tt), nil
	}

	return Time(time.Now()), fmt.Errorf("unable to parse date: %s", date)
}

// TryParseTime tries to parse time using a couple of formats before giving up
func TryParseTime(value string) (time.Time, bool) {
	var t time.Time
	var err error
	for _, layout := range formats {
		t, err = time.Parse(layout, value)
		if err == nil {
			return t, true
		}
	}

	return t, false
}
