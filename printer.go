package main

import (
	"fmt"
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

func Print(config Config, query Query, response QueryResponse) {
	for _, result := range response.FileResults {
		fmt.Printf("%s\n",
			strings.TrimSpace(colorize(result.Path, result.Bounds)),
		)

	}

	if len(response.FileResults) > 0 {
		fmt.Printf("\n")
	}

	for _, result := range response.Results {
		lineNoStr := strconv.Itoa(result.Lno)
		fmt.Printf("%s:%s:%s\n",
			color.CyanString(result.Path),
			color.YellowString(lineNoStr),
			strings.TrimSpace(colorize(result.Line, result.Bounds)),
		)
	}

}
