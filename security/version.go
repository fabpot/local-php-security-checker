package security

import (
	"encoding/json"
	"regexp"
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

	// be more lenient with pre-release versions, convert "2.0.0-alpha12" to "2.0.0-alpha.12"
	re := regexp.MustCompile(`(?i)(alpha|beta|rc)(\d+)$`)
	tmp = re.ReplaceAllString(tmp, "$1.$2")

	*v = Version(tmp)
	return nil
}
