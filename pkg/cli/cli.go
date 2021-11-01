package cli

type CommandUsage struct {
	Name string
	Args map[string]string	// Taken arguments and their meaning.
	Usage string			// Short description of command.
}

type Command interface {
	Run([]string) error
	Usage() CommandUsage
	Description() string
}

