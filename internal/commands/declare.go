package commands

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/deltron-fr/gash/internal/shell"
)

var (
	ErrDeclareMissingVariable = errors.New("missing variable")
	ErrDeclareInvalidFormat   = errors.New("declare: invalid format")
	ErrInvalidIdentifier      = errors.New("invalid identifier")
)

func declareArgs() map[string]declareOptions {
	options := map[string]declareOptions{
		"-p": {
			Name:        "-p",
			Description: "prints a description of the variable name",
		},
	}
	return options
}

func Declare(sh *shell.Shell, cmd *shell.Command) error {
	options := declareArgs()

	switch len(cmd.Args) {
	case 1:
		parts := strings.Split(cmd.Args[0], "=")

		if len(parts) == 1 {
			fmt.Fprintln(cmd.Stderr, ErrDeclareInvalidFormat.Error())
			return ErrDeclareInvalidFormat
		}

		if !validateVarName(parts[0]) {
			fmt.Fprintf(cmd.Stderr, "declare: `%s': not a valid identifier\n", cmd.Args[0])
			return ErrInvalidIdentifier
		}

		sh.EnvVariables[parts[0]] = parts[1]
	case 2:
		if _, exists := options[cmd.Args[0]]; !exists {
			fmt.Fprintf(cmd.Stderr, "%s: invalid option\n", cmd.Args[0])
			return ErrInvalidOptions
		} else {
			varName := cmd.Args[1]
			value, ok := sh.EnvVariables[varName]

			if !ok {
				fmt.Fprintf(cmd.Stderr, "declare: %s: not found\n", varName)
				return ErrDeclareMissingVariable
			}

			fmt.Fprintf(cmd.Stdout, "declare -- %s=\"%s\"\n", varName, value)
		}
	}

	return nil
}

func validateVarName(name string) bool {
	if len(name) == 0 {
		return false
	}

	if !isASCIILetter(name[0]) && name[0] != '_' {
		return false
	}

	for _, c := range name {
		if !unicode.IsDigit(c) && !unicode.IsLetter(c) && c != '_' {
			return false
		}
	}

	return true
}

func isASCIILetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

type declareOptions struct {
	Name        string
	Description string
}
