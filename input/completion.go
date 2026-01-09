package input

import (
	"os"
	"strings"
	"fmt"

	"github.com/deltron-fr/dshell/commands"
)

func autoCompletion(input string) []byte {
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