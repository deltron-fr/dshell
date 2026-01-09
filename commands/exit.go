package commands

import "os"

func handleExit(cmdName, redirection string, args ...string) {
	os.Exit(0)
}
