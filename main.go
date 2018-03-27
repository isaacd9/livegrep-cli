package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
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

	for _, result := range response.Results {
		fmt.Printf("%s:%d:# %s\n",
			result.Path, result.Lno, strings.TrimSpace(colorize(result.Line, result.Bounds)))
	}
}

func main() {
	flag.Parse()

	query := flag.Args()[0]
	runQuery(query)
}
