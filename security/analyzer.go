package security

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
)

// Vulnerabilities stores vulnerabilities for a lock file
type Vulnerabilities map[string]Vulnerability

// Vulnerability represents an vulnerability
type Vulnerability struct {
	Version    string           `json:"version"`
	Advisories []SimpleAdvisory `json:"advisories"`
}

// SimpleAdvisory represents an advisory for export
type SimpleAdvisory struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	CVE   string `json:"cve"`
}

func (a SimpleAdvisory) String() string {
	str := a.Title
	if a.CVE != "" {
		str = a.CVE + ": " + str
	}
	if a.Link != "" {
		str = str + " - " + a.Link
	}
	return str
}

// CountVulnerablePackages returns the number of packages with vulnerabilities
func (v *Vulnerabilities) CountVulnerablePackages() int {
	return len(*v)
}

// Count returns the number of vulnerabilities
func (v *Vulnerabilities) Count() int {
	count := 0
	for _, vs := range *v {
		count += len(vs.Advisories)
	}
	return count
}

// Keys returns package names in alpha order
func (v *Vulnerabilities) Keys() []string {
	keys := make([]string, len(*v))
	i := 0
	for k := range *v {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// Get returns the list of vulnerabilities for a given package
func (v *Vulnerabilities) Get(pkg string) *Vulnerability {
	vuln, ok := map[string]Vulnerability(*v)[pkg]
	if !ok {
		return nil
	}
	return &vuln
}

// Analyze checks if a give lock references packages with known security issues
func Analyze(lock *Lock, db *AdvisoryDB, noDevPackages bool) *Vulnerabilities {
	vulnerabilities := make(Vulnerabilities)
	packages := lock.Packages
	if !noDevPackages {
		packages = append(packages, lock.DevPackages...)
	}
	for _, p := range packages {
		var advs []SimpleAdvisory
		composerReference := "composer://" + p.Name
		packageBranchName := normalizeVersion(string(p.Version))
		for _, a := range db.Advisories {
			if a.Reference != composerReference {
				continue
			}
			for branchName, branch := range a.Branches {
				// dev versions must be checked via a date
				if isDev(p) {
					branchName = strings.TrimSuffix(branchName, ".x")
					if branchName != packageBranchName {
						continue
					}
					if time.Time(p.Time).IsZero() || time.Time(p.Time).After(time.Time(branch.Time)) {
						continue
					}
				} else {
					pv, err := version.NewVersion(string(p.Version))
					if err != nil {
						fmt.Fprintf(os.Stderr, "unable to parse version %s\n", p.Version)
						continue
					}
					constraintVersions := ""
					for _, v := range branch.Versions {
						constraintVersions = constraintVersions + ", " + string(v)
					}
					constraintVersions = strings.TrimPrefix(constraintVersions, ",")
					c, err := version.NewConstraint(constraintVersions)
					if err != nil {
						fmt.Fprintf(os.Stderr, "unable to parse version constraint %s\n", constraintVersions)
						continue
					}
					if !c.Check(pv) {
						continue
					}
				}
				advs = append(advs, SimpleAdvisory{
					CVE:   a.CVE,
					Link:  a.Link,
					Title: a.Title,
				})
			}
		}
		if len(advs) > 0 {
			vulnerabilities[p.Name] = Vulnerability{
				Version:    string(p.Version),
				Advisories: advs,
			}
		}
	}

	return &vulnerabilities
}

func normalizeVersion(version string) string {
	version = strings.TrimPrefix(version, "dev-")
	version = strings.TrimSuffix(version, "-dev")
	version = strings.TrimSuffix(version, ".x-dev")
	return version
}

// isDev checks if the package is a dev version
func isDev(p Package) bool {
	r := regexp.MustCompile("#.+$")
	version := r.ReplaceAllString(string(p.Version), "")
	if strings.HasPrefix(version, "dev-") || strings.HasSuffix(version, "-dev") {
		return true
	}

	return false
}
