package parser

type RedirectionCommands struct {
	Name        string
	Description string
}

func Redirection() map[string]RedirectionCommands {
	commands := map[string]RedirectionCommands{
		">": {
			Name:        ">",
			Description: "Redirect standard output",
		},
		"1>": {
			Name:        "1>",
			Description: "Redirect standard output",
		},
		"2>": {
			Name:        "2>",
			Description: "Redirect standard error",
		},
		">>": {
			Name:        ">>",
			Description: "Appending redirect standard output",
		},
		"1>>": {
			Name:        "1>>",
			Description: "Appending redirect standard output",
		},
		"2>>": {
			Name:        "2>>",
			Description: "Appending redirect standard error",
		},
	}
	return commands

}
