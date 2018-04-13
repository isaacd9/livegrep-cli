package main

import (
	"fmt"
	"github.com/fatih/color"
	"strconv"
	"strings"
)

var bgYellow = color.New(color.BgYellow).SprintFunc()
var fgYellow = color.New(color.FgYellow).SprintFunc()
var cyan = color.New(color.FgCyan, color.Bold).SprintFunc()

func colorize(str string, bounds []int) string {
	fst := str[0:bounds[0]]
	match := str[bounds[0]:bounds[1]]
	rest := str[bounds[1]:]
	return fmt.Sprintf("%s%s%s", fst, bgYellow(match), rest)
}

func makePrefix(config Config, result QueryResult, lineNoAdjust int) string {
	lineNoStr := strconv.Itoa(result.Lno + lineNoAdjust)
	var prefix string
	var filePrefix string
	var treePrefix string

	if config.printTree {
		treePrefix = cyan(result.Tree) + " "
	} else {
		filePrefix = ""
	}

	if !config.noPrintHeaders {
		filePrefix = fmt.Sprintf(
			"%s%s:",
			treePrefix,
			cyan(result.Path),
		)
	} else {
		filePrefix = ""
	}

	if !config.noPrintLineNumber {
		prefix = fmt.Sprintf("%s%s:",
			filePrefix,
			fgYellow(lineNoStr),
		)
	} else {
		prefix = filePrefix
	}
	return prefix
}

func buildLine(config Config, result QueryResult) string {
	return fmt.Sprintf("%s%s\n",
		makePrefix(config, result, 0),
		colorize(result.Line, result.Bounds),
	)
}

// Print prints baesed on a CLI config
func Print(config Config, query Query, response QueryResponse) {
	color.NoColor = !config.colorize
	lineCount := 0

	if !config.findInBody || config.findInFilename {
		for _, result := range response.FileResults {
			if config.numLines >= 0 && lineCount >= config.numLines {
				return
			}

			fmt.Printf("%s\n",
				strings.TrimSpace(colorize(result.Path, result.Bounds)),
			)
			lineCount++
		}
	}

	if !config.findInFilename || config.findInBody {
		for _, result := range response.Results {
			if config.numLines >= 0 && lineCount >= config.numLines {
				return
			}

			start := len(result.ContextBefore) - config.contextLinesBefore
			if start < 0 {
				start = 0
			}
			for i := start; i < len(result.ContextBefore); i++ {
				fmt.Printf(
					"%s%s\n",
					makePrefix(config, result, -(len(result.ContextBefore)-i)),
					result.ContextBefore[i],
				)
			}

			fmt.Print(buildLine(config, result))
			lineCount++

			for i := 0; i < config.contextLinesAfter; i++ {
				if i >= len(result.ContextAfter) {
					break
				}
				fmt.Printf(
					"%s%s\n",
					makePrefix(config, result, i+1),
					result.ContextAfter[i],
				)
			}
		}
	}
}
