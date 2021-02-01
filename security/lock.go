package security

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Package represents a Composer package
type Package struct {
	Name    string  `json:"name"`
	Version Version `json:"version"`
	Time    Time    `json:"time,omitempty"`
}

// Lock represents a Composer lock file
type Lock struct {
	Packages    []Package `json:"packages"`
	DevPackages []Package `json:"packages-dev"`
}

// NewLock creates a lock file wrapper
func NewLock(reader io.Reader) (*Lock, error) {
	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.New("unable to read lock file")
	}
	var l *Lock
	if err = json.Unmarshal(contents, &l); err != nil {
		return nil, errors.New("lock file is not valid JSON (not a composer.lock file?)")
	}
	if l.Packages == nil && l.DevPackages == nil {
		return nil, errors.New("lock file is not valid (no packages and no dev packages)")
	}
	return l, nil
}

// LocateLock locates a composer.lock
func LocateLock(path string) (io.Reader, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return nil, err
	}

	if path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(cwd, "composer.lock")
	} else if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		path = filepath.Join(path, "composer.lock")
	} else if strings.HasSuffix(path, "composer.json") {
		path = strings.Replace(path, "composer.json", "composer.lock", 1)
	}

	reader, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid lock file: %s", path, err)
	}

	return reader, nil
}
