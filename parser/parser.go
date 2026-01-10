package parser

import (
	"fmt"
	"strings"
)

func ParseInput(input string) []string {
	var args []string
	var current strings.Builder

	inQuotes := false
	inBackSlash := false
	var quote rune

	runes := []rune(input)

	for i, r := range runes {

		switch {
		case inBackSlash:
			current.WriteRune(r)
			inBackSlash = false

		case r == '\\' && inQuotes && quote == '\'':
			current.WriteRune(r)

		case r == '\\' && inQuotes && quote == '"':
			nextIndex := i + 1
			if nextIndex >= len(runes) {
				fmt.Println("malformed command input")
				return nil
			}
			if runes[nextIndex] == '\\' || runes[nextIndex] == '"' {
				inBackSlash = true
			} else {
				current.WriteRune(r)
			}

		case r == '\\' && !inQuotes:
			inBackSlash = true

		case r == '"' && !inBackSlash && (!inQuotes || quote == '"'):
			if !inQuotes {
				inQuotes = true
				quote = '"'
			} else {
				inQuotes = false
				quote = 0
			}

		case r == '\'' && !inBackSlash && (!inQuotes || quote == '\''):
			if !inQuotes {
				inQuotes = true
				quote = '\''
			} else {
				inQuotes = false
				quote = 0
			}

		case r == ' ' && !inQuotes:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}

		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}
