package main

import (
	"fmt"
)

func handleType(args ...string) {
	availableCmds := Commands()
	for _, arg := range args {
		fmt.Println(arg)
		v, exists := availableCmds[arg]
		fmt.Println(v.name)
		if exists {
			fmt.Printf("%s is a shell builtin\n", arg)
		} else {
			checkPath(arg, "type")
		}
	}
}

