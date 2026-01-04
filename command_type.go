package main

import (
	"fmt"
	"os"
)

func handleType(cmdName, redirection string, args ...string) {
	availableCmds := Commands()

	if redirection == "" {
		for _, arg := range args {
			_, exists := availableCmds[arg]
			if exists {
				fmt.Printf("%s is a shell builtin\n", arg)
			} else {
				checkPath(nil, arg, "type")
			}
		}
		return
	}

	filepath := args[len(args)-1]
	args = args[:len(args)-2]

	switch redirection {
	case ">", "1>":
		file, err := os.Create(filepath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		for _, arg := range args {
			_, exists := availableCmds[arg]
			if exists {
				fmt.Fprintf(file, "%s is a shell builtin\n", arg)
				fmt.Fprintf(file, "\n")
			} else {
				checkPath(file, arg, "type")
				}
			}
		case "2>":
			file, err := os.Create(filepath)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			for _, arg := range args {
				_, exists := availableCmds[arg]
				if exists {
					_, err = fmt.Fprintf(file, "%s is a shell builtin\n", arg)
					if err != nil {
						fmt.Fprintln(file, err)
						return
					}
					fmt.Fprintf(file, "\n")
				} else {
					checkPath(file, arg, "type")
					}
				}
		}
}

