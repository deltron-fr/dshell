package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type commandFunc func(...string)

type builtInCommands struct {
	name        string
	description string
	callback     commandFunc
}

func startRepl() {
	for {
		fmt.Print("$ ")

		var input string
		scanner := bufio.NewScanner(os.Stdin)

		if scanner.Scan() {
			input = scanner.Text()
		}

		if input == "" {
			continue
		}

		var cmd string
		var extraArgs []string
		args := strings.Fields(input)

		cmd = args[0]
		if len(args) > 1 {
			extraArgs = args[1:]
		}

		commands := Commands()
		if command, exists := commands[cmd]; exists {
			command.callback(extraArgs...)
		} else {
			handleUnknownCommand(cmd)
		}
	}
}

func Commands() map[string]builtInCommands {
	commands := map[string]builtInCommands{
		"exit": {
			name: "exit",
			description: "Exit the shell",
			callback: handleExit,
		},
		"echo": {
			name: "echo",
			description: "display a line of text",
			callback: handleEcho,
		},
		"type": {
			name: "type",
			description: "display information about command type",
			callback: handleType,
		},
	}

	return commands
}

func handleUnknownCommand(cmd string) {
	fmt.Printf("%s: command not found\n", cmd)
}


