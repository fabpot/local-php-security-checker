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

func printHeader() {
	fmt.Println(" _                    _   ____  _   _ ____  ")
	fmt.Println("| |    ___   ___ __ _| | |  _ \\| | | |  _ \\	")
	fmt.Println("| |   / _ \\ / __/ _` | | | |_) | |_| | |_) |")
	fmt.Println("| |__| (_) | (_| (_| | | |  __/|  _  |  __/ ")
	fmt.Println("|_____\\___/ \\___\\__,_|_| |_|   |_| |_|_|    ")
	fmt.Println("")
	fmt.Println(" ____                       _ _            ____ _               _ ")
	fmt.Println("/ ___|  ___  ___ _   _ _ __(_) |_ _   _   / ___| |__   ___  ___| | _____ _ __ ")
	fmt.Println("\\___ \\ / _ \\/ __| | | | '__| | __| | | | | |   | '_ \\ / _ \\/ __| |/ / _ \\ '__|")
	fmt.Println(" ___) |  __/ (__| |_| | |  | | |_| |_| | | |___| | | |  __/ (__|   <  __/ |   ")
	fmt.Println("|____/ \\___|\\___|\\__,_|_|  |_|\\__|\\__, |  \\____|_| |_|\\___|\\___|_|\\_\\___|_|   ")
	fmt.Println("                                  |___/                                       ")
	fmt.Printf("%s, built at %s\n", version, date)
}

func main() {
	format := flag.String("format", "ansi", "Output format (ansi, markdown, json, or yaml)")
	path := flag.String("path", "", "composer.lock file or directory")
	local := flag.Bool("local", false, "Do not make HTTP calls (needs a valid cache file)")
	updateCacheOnly := flag.Bool("update-cache", false, "Update the cache (other flags are ignored)")
	help := flag.Bool("help", false, "Output help and version")
	quiet := flag.Bool("quiet", false, "Only print out vulnerable dependencies")
	ansi := flag.Bool("ansi", false, "no color")
	flag.Parse()

	if *ansi {
		security.Nocolor()
	}

	if !(*quiet) {
		security.Verbose = true
		printHeader()
	}

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	db, err := security.NewDB(*local)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to load the advisory DB: %s\n", err)
		os.Exit(127)
	}

	if *updateCacheOnly {
		if err := db.Load(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(127)
		}
		return
	}

	if *format != "" && *format != "markdown" && *format != "json" && *format != "yaml" && *format != "ansi" {
		fmt.Fprintf(os.Stderr, "format \"%s\" is not supported (supported formats: markdown, ansi, json, and yaml)\n", *format)
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

	vulns := security.Analyze(lock, db)

	output, err := security.Format(vulns, *format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to output the results: %s\n", err)
		os.Exit(127)
	}
	fmt.Printf(string(output))

	if vulns.Count() > 0 {
		os.Exit(1)
	}
}
