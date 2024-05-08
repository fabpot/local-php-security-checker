package security

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Format formats the vulnerabilities in the given format
func Format(vulns *Vulnerabilities, format string) ([]byte, error) {
	if format == "ansi" {
		return ToANSI(vulns), nil
	} else if format == "text" || format == "txt" || format == "markdown" || format == "md" {
		return ToMarkdown(vulns), nil
	} else if format == "json" {
		return ToJSON(vulns, true)
	} else if format == "raw_json" {
		return ToJSON(vulns, false)
	} else if format == "junit" {
		return ToJunit(vulns)
	} else if format == "yaml" || format == "yml" {
		return ToYAML(vulns)
	}
	return nil, fmt.Errorf("unknown format %s", format)
}

// ToANSI returns vulnerabilities as text with ANSI code for colors
func ToANSI(vulns *Vulnerabilities) []byte {
	var output string
	output += "\u001B[33mSymfony Security Check Report\u001B[0m\n"
	output += "\u001B[33m=============================\u001B[0m\n\n"
	if vulns.CountVulnerablePackages() == 1 {
		output += "\u001B[41m1 package\u001B[0m has known vulnerabilities.\n"
	} else if vulns.CountVulnerablePackages() > 0 {
		output += fmt.Sprintf("\u001B[41m%d packages\u001B[0m have known vulnerabilities.\n", vulns.CountVulnerablePackages())
	} else {
		output += "\u001B[32mNo packages have known vulnerabilities.\u001B[0m"
	}
	output += fmt.Sprintln("")
	links := ""
	ref := 0
	for _, pkg := range vulns.Keys() {
		v := vulns.Get(pkg)
		str := fmt.Sprintf("%s (%s)", pkg, v.Version)
		output += fmt.Sprintf("\u001B[33m%s\u001B[0m\n\u001B[33m%s\u001B[0m\n\n", str, strings.Repeat("-", len(str)))
		for _, a := range v.Advisories {
			cve := a.CVE
			if cve == "" {
				ref++
				cve = fmt.Sprintf("CVE-NONE-%04d", ref)
			}
			title := strings.TrimPrefix(a.Title, a.CVE+": ")

			if a.Link == "" {
				output += fmt.Sprintf(" * \u001B[34m%s\u001B[0m: %s\n", cve, title)
			} else {
				output += fmt.Sprintf(" * [\u001B[34m%s\u001B[0m][]: %s\n", cve, title)
				links += fmt.Sprintf("[%s]: \u001B]8;;%s\u0007%s\u001B]8;;\u0007\u001B[0m\n", cve, a.Link, a.Link)
			}
		}
		output += fmt.Sprintln("")
	}
	output += links
	output += fmt.Sprintln("")

	output += "\u001B[33mNote that this checker can only detect vulnerabilities that are referenced in the security advisories database.\n" +
		"Execute this command regularly to check the newly discovered vulnerabilities.\u001B[0m\n"

	return []byte(output)
}

var ansiRe = regexp.MustCompile("(\u001B\\[\\d+m|\u001B\\]8;;.*?\u0007)")

// ToMarkdown returns vulnerabilities as Markdown
func ToMarkdown(vulns *Vulnerabilities) []byte {
	return ansiRe.ReplaceAll(ToANSI(vulns), nil)
}

// ToJSON outputs vulnerabilities as JSON
func ToJSON(vulns *Vulnerabilities, prettify bool) ([]byte, error) {
	if prettify {
		return json.MarshalIndent(vulns, "", "    ")
	}

	return json.Marshal(vulns)
}

// ToYAML outputs vulnerabilities as YAML
func ToYAML(vulns *Vulnerabilities) ([]byte, error) {
	return yaml.Marshal(vulns)
}
