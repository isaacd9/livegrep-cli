package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
)

// Config represents the configuration of the CLI tool
type Config struct {
	caseInsensitive    bool
	colorize           bool
	contextLinesAfter  int
	contextLinesBefore int
	fixedStrings       bool
	findInFilename     bool
	findInBody         bool
	noPrintHeaders     bool
	numLines           int
	pattern            string
	printFilename      bool
	printLineNumber    bool
}

func initFlags() Config {
	c := &Config{}
	var contextLines int
	flag.IntVar(
		&contextLines,
		"C",
		0,
		"Print num lines of leading and trailing context surrounding each match.",
	)
	flag.IntVar(
		&c.contextLinesAfter,
		"A",
		0,
		"Print num lines of trailing context after each match.",
	)
	flag.IntVar(
		&c.contextLinesBefore,
		"B",
		0,
		"Print num lines of trailing context before each match.",
	)
	flag.IntVar(
		&c.numLines,
		"c",
		-1,
		"Only a count of selected lines is written to standard output.",
	)
	flag.BoolVar(
		&c.colorize,
		"color",
		true,
		"Mark up the matching text in color",
	)
	flag.StringVar(
		&c.pattern,
		"e",
		"",
		"Specify a pattern used during the search of the input: an input line is"+
			"selected if it matches any of the specified patterns.",
	)
	flag.BoolVar(
		&c.fixedStrings,
		"F",
		false,
		"Interpret pattern as a set of fixed strings",
	)
	flag.BoolVar(
		&c.noPrintHeaders,
		"h",
		false,
		"Never print filename headers (i.e. filenames) with output lines.",
	)
	flag.BoolVar(
		&c.caseInsensitive,
		"i",
		false,
		"Perform case insensitive matching.",
	)
	flag.BoolVar(
		&c.printLineNumber,
		"n",
		true,
		"Each output line is preceded by its relative line number in the file.",
	)
	flag.BoolVar(
		&c.findInFilename,
		"f",
		true,
		"Look in the names of files for matches (like find)",
	)
	flag.BoolVar(
		&c.findInBody,
		"b",
		true,
		"Look in the contents of files for matches (like grep)",
	)
	var version bool
	flag.BoolVar(&version, "v", false, "Print version and exit")

	flag.Parse()

	if version {
		fmt.Printf("0.0.1\n")
		os.Exit(0)
	}

	if contextLines > c.contextLinesAfter {
		c.contextLinesAfter = contextLines
	}
	if contextLines > c.contextLinesBefore {
		c.contextLinesBefore = contextLines
	}

	return *c
}

func main() {
	config := initFlags()

	if len(flag.Args()) == 0 {
		os.Exit(0)
	}

	query := flag.Args()[0]
	var host string
	if os.Getenv("LIVEGREP_HOST") != "" {
		host = os.Getenv("LIVEGREP_HOST")
	} else {
		host = "livegrep.com"
	}

	l := NewLivegrep(host)
	if os.Getenv("LIVEGREP_USE_HTTPS") != "" {
		l.UseHTTPS = true
	}

	unixSocket := os.Getenv("LIVEGREP_UNIX_SOCKET")
	if unixSocket != "" {
		transport := &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", unixSocket)
			},
		}

		httpClient := &http.Client{
			Transport: transport,
		}

		l.Client = httpClient
	}

	q := l.NewQuery(query)

	q.FoldCase = config.caseInsensitive
	q.Regex = !config.fixedStrings

	if len(config.pattern) > 0 {
		q.Term = config.pattern
	}

	response, err := l.Query(q)
	if err != nil {
		panic(err)
	}

	Print(config, q, response)
}
