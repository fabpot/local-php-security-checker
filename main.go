package main

/*

	Checks security issues in project dependencies. Without arguments, it looks
	for a "composer.lock" file in the current directory. Pass it explicitly to check
	a specific "composer.lock" file.

*/

import (
	"flag"
	"fmt"
	"os"

	"github.com/fabpot/local-php-security-checker/security"
)

var (
	version = "dev"
	date    = "unknown"
)

func main() {
	format := flag.String("format", "ansi", "Output format (ansi, junit, markdown, json, or yaml)")
	path := flag.String("path", "", "composer.lock file or directory")
	advisoryArchiveURL := flag.String("archive", security.AdvisoryArchiveURL, "Advisory archive URL")
	cacheDir := flag.String("cache-dir", os.TempDir(), "Cache directory")
	local := flag.Bool("local", false, "Do not make HTTP calls (needs a valid cache file)")
	noDevPackages := flag.Bool("no-dev", false, "Do not check packages listed under require-dev")
	updateCacheOnly := flag.Bool("update-cache", false, "Update the cache (other flags are ignored)")
	help := flag.Bool("help", false, "Output help and version")
	flag.Parse()

	if *help {
		fmt.Printf("Local PHP Security Checker %s, built at %s\n", version, date)
		flag.Usage()
		os.Exit(0)
	}

	db, err := security.NewDB(*local, *advisoryArchiveURL, *cacheDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to load the advisory DB: %s\n", err)
		os.Exit(127)
	}

	if *updateCacheOnly {
		if err := db.Load(*advisoryArchiveURL); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(127)
		}
		return
	}

	if *format != "" && *format != "markdown" && *format != "json" && *format != "yaml" && *format != "ansi" && *format != "junit" {
		fmt.Fprintf(os.Stderr, "format \"%s\" is not supported (supported formats: markdown, ansi, json, junit, and yaml)\n", *format)
		os.Exit(2)
	}

	lockReader, err := security.LocateLock(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(127)
	}

	lock, err := security.NewLock(lockReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to load the lock file: %s\n", err)
		os.Exit(127)
	}

	vulns := security.Analyze(lock, db, *noDevPackages)

	output, err := security.Format(vulns, *format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to output the results: %s\n", err)
		os.Exit(127)
	}
	fmt.Print(string(output))

	if vulns.Count() > 0 {
		os.Exit(1)
	}
}
