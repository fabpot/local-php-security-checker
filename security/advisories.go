package security

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// AdvisoryArchiveURL represents the advisories database URL
const AdvisoryArchiveURL = "https://codeload.github.com/FriendsOfPHP/security-advisories/zip/master"

// AdvisoryDB stores all known security advisories
type AdvisoryDB struct {
	Advisories  []Advisory
	cacheDir    string
	noHTTPCalls bool
}

// Advisory represents a single security advisory
type Advisory struct {
	Title     string             `yaml:"title"`
	Link      string             `yaml:"link"`
	CVE       string             `yaml:"cve"`
	Branches  map[string]*Branch `yaml:"branches"`
	Reference string             `yaml:"reference"`
}

// Branch represents a Git branch
type Branch struct {
	Versions []string `yaml:"versions"`
	Time     Time     `yaml:"time"`
}

// Cache stores the Github response to save bandwith
type Cache struct {
	Key  string
	Date string
	Body []byte
}

// NewDB fetches the advisory DB from Github
func NewDB(noHTTPCalls bool, advisoryArchiveURL, cacheDir string) (*AdvisoryDB, error) {
	db := &AdvisoryDB{noHTTPCalls: noHTTPCalls, cacheDir: cacheDir}
	if err := db.load(advisoryArchiveURL); err != nil {
		return nil, fmt.Errorf("unable to fetch advisories: %s", err)
	}

	return db, nil
}

// load loads fetches the database from Github and reads/loads current advisories
// from the repository. Cache handling is delegated to http.Transport and
// **must** be handled appropriately.
func (db *AdvisoryDB) load(advisoryArchiveURL string) error {
	if len(db.Advisories) > 0 {
		return nil
	}

	db.Advisories = []Advisory{}

	var cache *Cache
	cachePath := filepath.Join(db.cacheDir, "php_sec_db.json")
	if cacheContent, err := ioutil.ReadFile(cachePath); err == nil {
		// ignore errors
		json.Unmarshal(cacheContent, &cache)
	}

	if db.noHTTPCalls && cache == nil {
		return errors.New("--local can only be used when a local HTTP cache is available")
	}

	if !db.noHTTPCalls {
		req, err := http.NewRequest("GET", advisoryArchiveURL, nil)
		if err != nil {
			return err
		}
		if cache != nil {
			req.Header.Add("If-None-Match", cache.Key)
			req.Header.Add("If-Modified-Since", cache.Date)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		var body []byte
		if resp.StatusCode != http.StatusNotModified {
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			key := resp.Header.Get("ETag")
			date := resp.Header.Get("Date")
			if key != "" || date != "" {
				cache = &Cache{Key: key, Date: date, Body: body}
			}
			cacheContent, err := json.Marshal(cache)
			if err == nil {
				ioutil.WriteFile(cachePath, cacheContent, 0644)
			}
		}
	}

	zipReader, err := zip.NewReader(bytes.NewReader(cache.Body), int64(len(cache.Body)))
	if err != nil {
		return err
	}

	// Read all the files from the zip archive
	for _, zipFile := range zipReader.File {
		if !strings.HasSuffix(zipFile.Name, ".yaml") {
			continue
		}
		f, err := zipFile.Open()
		if err != nil {
			return err
		}
		defer f.Close()

		contents, err := ioutil.ReadAll(f)
		if err != nil {
			return fmt.Errorf("unable to read %s: %s", zipFile.Name, err)
		}

		var pa Advisory
		if err := yaml.Unmarshal(contents, &pa); err != nil {
			return fmt.Errorf("%s is not a valid YAML file: %s", zipFile.Name, err)
		}

		db.Advisories = append(db.Advisories, pa)
	}

	return nil
}
