package main

import (
	"flag"
	"fmt"
	//"github.com/chzyer/readline"
	"github.com/fatih/color"
	"strconv"
	"strings"
)

func colorize(str string, bounds []int) string {
	fst := str[0:bounds[0]]
	match := str[bounds[0]:bounds[1]]
	rest := str[bounds[1]:]
	return fmt.Sprintf("%s%s%s", fst, color.YellowString(match), rest)
}

func runQuery(query string) {
	l := &Livegrep{URL: "livegrep.com"}
	q := l.NewQuery(query)
	response, err := l.Query(q)
	if err != nil {
		panic(err)
	}

	for _, result := range response.FileResults {
		fmt.Printf("%s\n",
			strings.TrimSpace(colorize(result.Path, result.Bounds)),
		)
	}

	if len(response.FileResults) > 0 {
		fmt.Printf("\n")
	}

	for _, result := range response.Results {
		fmt.Printf("%s:%s:%s\n",
			color.CyanString(result.Path),
			color.YellowString(strconv.Itoa(result.Lno)),
			strings.TrimSpace(colorize(result.Line, result.Bounds)),
		)
	}
}

func shell() {
}

func main() {
	flag.Parse()

	if len(flag.Args()) > 0 {
		query := flag.Args()[0]
		runQuery(query)
	} else {
		shell()
	}
}
