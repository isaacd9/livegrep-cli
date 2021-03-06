package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type repositoriesList []string

func (r *repositoriesList) String() string {
	return strings.Join(*r, ",")
}

func (r *repositoriesList) Set(value string) error {
	*r = append(*r, value)
	return nil
}

// Config represents the configuration of the CLI tool
type Config struct {
	caseInsensitive    bool
	colorize           bool
	contextLinesAfter  int
	contextLinesBefore int
	findInBody         bool
	findInFilename     bool
	fixedStrings       bool
	noPrintHeaders     bool
	noPrintLineNumber  bool
	numLines           int
	pattern            string
	printFilename      bool
	printTree          bool
	repositories       repositoriesList
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
		&c.noPrintLineNumber,
		"N",
		false,
		"Do not precede each output line is by its relative line number in the file.",
	)
	flag.BoolVar(
		&c.findInFilename,
		"f",
		false,
		"Look in the names of files for matches (like find)",
	)
	flag.BoolVar(
		&c.findInBody,
		"b",
		false,
		"Look in the contents of files for matches (like grep)",
	)
	flag.BoolVar(
		&c.printTree,
		"t",
		false,
		"Print tree in addition to file path",
	)
	flag.Var(
		&c.repositories,
		"r",
		"Repositories to search in. Specify -r multiple times to specify a list of repositories.",
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
	if os.Getenv("LIVEGREP_USE_HTTPS") == "false" {
		l.UseHTTPS = false
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
	q.Repositories = config.repositories

	q.FoldCase = config.caseInsensitive
	q.Regex = !config.fixedStrings

	if len(config.pattern) > 0 {
		q.Term = config.pattern
	}

	response, err := l.Query(q)
	if err != nil {
		log.Fatalln(err)
	}

	Print(config, q, response)
}
