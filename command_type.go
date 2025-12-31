package main

import (
	"fmt"
)

func handleType(cmdName, redirection string, args ...string) {
	availableCmds := Commands()
	for _, arg := range args {
		_, exists := availableCmds[arg]
		if exists {
			fmt.Printf("%s is a shell builtin\n", arg)
		} else {
			checkPath(arg, "type")
		}
	}
}
