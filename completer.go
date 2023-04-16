package main

import (
	"strings"

	"github.com/peterh/liner"
)

func newWordCompleter() liner.WordCompleter {
	return func(line string, pos int) (head string, completions []string, tail string) {
		head, partial, tail := splitLine(line, pos)

		if partial == "" {
			// match to nothing
			return
		}

		// function auto-completes
		for name := range funcMap() {
			if strings.HasPrefix(name, partial) {
				completions = append(completions, name)
			}
		}
		return
	}
}

func splitLine(line string, pos int) (head, partial, tail string) {
	startPos := partialStartPos(line, pos, "{}()<>=| ")
	endPos := partialEndPos(line, pos, "{}()<>=| ")

	head = subStr(line, 0, startPos)
	tail = subStr(line, endPos, len(line))
	partial = subStr(line, startPos, endPos)
	return
}

func partialStartPos(line string, pos int, delims string) int {
	delimPos := strings.LastIndexAny(subStr(line, 0, pos), delims)
	if delimPos < 0 {
		return 0
	}
	return delimPos + 1
}

func partialEndPos(line string, pos int, delims string) int {
	relativeDelimPos := strings.IndexAny(subStr(line, pos, len(line)), delims)
	if relativeDelimPos < 0 {
		return len(line)
	}
	return pos + relativeDelimPos
}

func subStr(str string, start, end int) string {
	if start >= len(str) {
		return ""
	}
	if end <= 0 {
		return ""
	}

	if end > len(str) {
		end = len(str)
	}
	if start < 0 {
		start = 0
	}

	return str[start:end]
}
