package security

import (
	"encoding/json"
)

// Version represents a composer.json version (can be a string or an integer)
type Version string

// UnmarshalJSON converts versions as integers to strings
func (v *Version) UnmarshalJSON(b []byte) error {
	var tmpNumber json.Number
	if err := json.Unmarshal(b, &tmpNumber); err == nil {
		if _, err := tmpNumber.Int64(); err == nil {
			*v = Version(tmpNumber.String())
			return nil
		}
	}
	var tmp string
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*v = Version(tmp)
	return nil
}
