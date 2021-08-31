package security

import (
	"encoding/xml"
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
	Failure   string   `xml:"failure,omitempty"`
}
