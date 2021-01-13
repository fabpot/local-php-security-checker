package security

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeParsing(t *testing.T) {
	date := time.Date(2018, 11, 22, 23, 7, 0, 0, time.UTC)
	for _, testCase := range []string{
		"2018-11-22 23:07 UTC",
		"2018-11-22 23:07:00",
		"2018-11-22 23:07:00 UTC",
	} {
		tt, ok := TryParseTime(testCase)
		assert.True(t, ok)
		assert.NotNil(t, tt)
		assert.Equal(t, date.Format(time.RFC822Z), tt.Format(time.RFC822Z))
	}

	date = time.Date(2018, 11, 22, 23, 7, 0, 0, time.UTC)
	for _, testCase := range []string{
		"2018-11-22 23:07:00 +00:00",
		"2018-11-22 23:07 +00:00",
	} {
		tt, ok := TryParseTime(testCase)
		assert.True(t, ok)
		assert.NotNil(t, tt)
		assert.Equal(t, date.Format(time.RFC822Z), tt.Format(time.RFC822Z))
	}
}
