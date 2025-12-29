package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func handleType(args ...string) {
	availableCmds := Commands()
	for _, arg := range args {
		_, exists := availableCmds[arg]
		if exists {
			fmt.Printf("%s is a shell builtin\n", arg)
		} else {
			checkPath(arg)
		}
	}
}

func checkPath(cmdName string) {
	pathEnv := os.Getenv("PATH")
	separator := string(os.PathListSeparator)

	directories := strings.Split(pathEnv, separator)
	for _, dir := range directories {
		cmdPath := dir + "/" + cmdName
		if !fileExists(cmdPath) {
			continue
		}

		if !isExecutable(cmdPath) {
			continue
		}

		fmt.Printf("%s is %s\n", cmdName, cmdPath)
		return
	}
	fmt.Printf("%s: not found\n", cmdName)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	return err == nil
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		log.Printf("error getting file info: %v", err)
		return false
	}

	mode := info.Mode()
	isExec := mode&0100 != 0
	return isExec
}