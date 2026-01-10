package input

import (
	"fmt"
	"os"
	"strings"

	"github.com/deltron-fr/dshell/commands"
)

func autoCompletion(input string) []byte {
	out := autoCompleteCmds(input)
	if len(out) != 0 {
		return out
	}

	out = autoCompleteFiles(input)
	if len(out) != 0 {
		return out
	}

	out = autoCompleteCmdPath(input)
	if len(out) != 0 {
		return out
	}

	return []byte{}
}

func autoCompleteCmds(input string) []byte {
	var cmdName string

	commands := commands.Commands()
	for _, v := range commands {
		if strings.HasPrefix(v.Name, input) {
			cmdName = v.Name
			break
		}
	}

	if cmdName == "" {
		return []byte{}
	}

	return []byte(cmdName[len(input):])
}

func autoCompleteCmdPath(input string) []byte {
	pathEnv := os.Getenv("PATH")
	separator := string(os.PathListSeparator)

	directories := strings.Split(pathEnv, separator)
	var extCmd string

	for _, dir := range directories {
		files, err := os.ReadDir(dir)
		if err == os.ErrPermission {
			fmt.Fprintf(os.Stderr, "insufficient permission to read directory: %v", dir)
			continue
		}

		for _, f := range files {
			if !f.Type().IsRegular() {
				continue
			}

			if strings.HasPrefix(f.Name(), input) {
				extCmd = f.Name()
				break
			}
		}
		if extCmd != "" {
			break
		}
	}

	if extCmd == "" {
		return []byte{}
	}

	return []byte(extCmd[len(input):])
}

func autoCompleteFiles(input string) []byte {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting current working directory: %v", err)
		return []byte{}
	}

	files, err := os.ReadDir(pwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading directory: %v", err)
		return []byte{}
	}

	var fileName string
	for _, f := range files {
		if strings.HasPrefix(f.Name(), input) {
			fileName = f.Name()
			break
		}
	}

	if fileName == "" {
		return []byte{}
	}

	return []byte(fileName[len(input):])
}
