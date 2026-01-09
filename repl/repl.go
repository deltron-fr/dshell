package repl

import (
	"fmt"
	
	"github.com/deltron-fr/dshell/input"
	"github.com/deltron-fr/dshell/parser"
	"github.com/deltron-fr/dshell/commands"
)

func StartRepl() {
	for {
		fmt.Print("$ ")

		input := input.RawModeHandler()
		if input == "" {
			continue
		}

		var cmd string
		var extraArgs []string
		args := parser.ParseInput(input)
		if args == nil {
			continue
		}

		cmd = args[0]
		if len(args) > 1 {
			extraArgs = args[1:]
		}

		invalid := false
		var redCmd parser.RedirectionCommands

		redirectCommands := parser.Redirection()
		for i, arg := range extraArgs {
			if c, ok := redirectCommands[arg]; ok {
				if i+1 >= len(extraArgs) {
					fmt.Println("invaid command input")
					invalid = true
					break
				} else {
					if i+1 != len(extraArgs)-1 {
						fmt.Println("invalid command input, too many destination arguments")
						break
					}
					redCmd = c
				}
			}
		}

		if invalid {
			continue
		}

		builtinCmds := commands.Commands()

		if command, exists := builtinCmds[cmd]; exists {
			command.Callback(command.Name, redCmd.Name, extraArgs...)
		} else {
			commands.HandleExec(cmd, redCmd.Name, extraArgs...)
		}
	}
}

