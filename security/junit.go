package security

import (
	"encoding/xml"
	"fmt"
)

type testsuites struct {
	XMLName    xml.Name `xml:"testsuites"`
	Name       string   `xml:"name,attr"`
	Testsuites []testsuite
}

type testsuite struct {
	XMLName   xml.Name `xml:"testsuite"`
	Package   string   `xml:"package,attr"`
	Errors    int      `xml:"errors,attr"`
	Failures  int      `xml:"failures,attr"`
	Tests     int      `xml:"tests,attr"`
	Testcases []testcase
}

type testcase struct {
	XMLName   xml.Name `xml:"testcase"`
	Name      string   `xml:"name,attr"`
	Classname string   `xml:"classname,attr"`
	Failure   []string `xml:"failure,omitempty"`
}

func ToJunit(vulns *Vulnerabilities) ([]byte, error) {
	var packages []testsuite
	var cases []testcase
	ts := testsuite{}
	for _, pkg := range vulns.Keys() {
		v := vulns.Get(pkg)
		tc := testcase{
			Classname: "packages",
			Name:      fmt.Sprintf("%s (%s)", pkg, v.Version),
		}
		for _, a := range v.Advisories {
			tc.Failure = append(tc.Failure, fmt.Sprintf("%s - %s (%s)", a.CVE, a.Title, a.Link))
		}
		cases = append(cases, tc)
		ts.Failures++
		ts.Tests++
	}
	ts.Testcases = cases
	packages = append(packages, ts)
	out := testsuites{
		Name:       "Symfony Security Check Report",
		Testsuites: packages,
	}
	return xml.MarshalIndent(&out, "  ", "    ")
}
