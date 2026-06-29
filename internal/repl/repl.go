package repl

import (
	"fmt"
	"os"
	"strings"

	"github.com/deltron-fr/gash/internal/commands"
	"github.com/deltron-fr/gash/internal/input"
	"github.com/deltron-fr/gash/internal/parser"
	"github.com/deltron-fr/gash/internal/shell"
)

func StartRepl() {
	// StartRepl runs the main read-eval-print loop. It prints a prompt,
	// reads a line using the raw-mode input handler, runs tab-completion
	// listing when requested, parses the input, checks for redirections,
	// and sends the command to builtins or external commands.
	//
	// `exit` builtin will terminate this process.
	var buffer string
	HistFile := os.Getenv("HISTFILE")

	sh := shell.NewShell()
	commands.RegisterBuiltins(sh)
	commands.LoadHistoryToMemory(sh, HistFile)

	for {
		sh.DrainJobUpdates()
		fmt.Print("$ ")

		input, tabMatches := input.RawModeHandler(*sh, buffer, sh.History)

		if len(tabMatches) > 0 {
			for _, match := range tabMatches {
				fmt.Fprintf(os.Stdout, "%s  ", match)
			}
			fmt.Println()
			buffer = input
			continue
		}

		buffer = ""

		if input == "" {
			continue
		}

		h := commands.AddEntry(input, sh.History)
		sh.History = append(sh.History, *h)

		args := parser.ParseInput(input)
		if args == nil {
			continue
		}

		pipeline, file, bg, err := ParsePipeline(sh, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			continue
		}

		sh.Executor(pipeline, bg)
		if file != nil {
			file.Close()
		}
	}
}

// ParsePipeline builds a pipeline from parsed args and applies any redirections.
// It returns the pipeline plus the last redirection file opened (if any).
func ParsePipeline(sh *shell.Shell, args []string) (*shell.Pipeline, *os.File, bool, error) {
	pipeline := shell.NewPipeline()
	isFirstArg := true
	isBackgroundJob := false
	var f *os.File

	cmd := shell.Command{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if args[len(args)-1] == "&" {
		isBackgroundJob = true
		args = args[:len(args)-1]
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "|":
			pipeline.Commands = append(pipeline.Commands, cmd)
			cmd = shell.Command{
				Stdin:  os.Stdin,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}
			isFirstArg = true

		case isFirstArg:
			fullParsedArg, err := handleParameterExpansion(sh, arg)
			if err != nil {
				return nil, nil, false, err
			}

			if fullParsedArg == "" {
				continue
			}

			cmd.Name = fullParsedArg
			isFirstArg = false

		case isRedirection(arg):
			if i+1 < len(args) {
				target := args[i+1]
				r := parser.Redirect{
					Operator: parser.Redirector(arg),
					Target:   target,
				}
				f, err := r.Apply(&cmd)
				if err != nil {
					return pipeline, f, isBackgroundJob, err
				}
				i++
			}

		default:
			fullParsedArg, err := handleParameterExpansion(sh, arg)
			if err != nil {
				return nil, nil, false, err
			}

			if fullParsedArg == "" {
				continue
			}

			cmd.Args = append(cmd.Args, fullParsedArg)
		}
	}

	if cmd.Name != "" {
		pipeline.Commands = append(pipeline.Commands, cmd)
	}

	return pipeline, f, isBackgroundJob, nil
}

// isRedirection reports whether a token is a supported redirection operator.
func isRedirection(token string) bool {
	redirectionOperators := parser.Redirection()
	if _, ok := redirectionOperators[token]; !ok {
		return false
	}

	return true
}

func handleParameterExpansion(sh *shell.Shell, arg string) (string, error) {
	length := len(arg)
	var fullParsedArg strings.Builder
	var fullIdx int

	for fullIdx < length {

		before, after, found := strings.Cut(arg[fullIdx:], "$")
		if !found {
			fullParsedArg.WriteString(before)
			break
		}

		fullParsedArg.WriteString(before)
		fullIdx += len(before)

		newArg, idx, err := parseVariable(sh, after)
		if err != nil {
			return "", err
		}

		fullIdx += idx
		fullParsedArg.WriteString(newArg)
	}

	return fullParsedArg.String(), nil
}

func parseVariable(sh *shell.Shell, arg string) (string, int, error) {
	var parsedWord strings.Builder
	var idx int

	isClosed := false

	if arg[0] != '{' {
		return resolveVariable(sh, arg), len(arg) + 1, nil
	}

	for _, s := range arg {
		idx++
		if s == '{' {
			continue
		}

		if s == '}' {
			isClosed = true
			break
		}

		parsedWord.WriteString(string(s))
	}

	if !isClosed {
		return "", 0, fmt.Errorf("variable is not closed")
	}

	return resolveVariable(sh, parsedWord.String()), idx + 1, nil
}

func resolveVariable(sh *shell.Shell, name string) string {
	v, ok := sh.EnvVariables[name]
	if !ok {
		name = ""
	} else {
		name = v
	}

	return name
}
