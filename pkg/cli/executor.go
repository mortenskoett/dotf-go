package cli

import (
	"fmt"

	"github.com/mortenskoett/dotf-go/pkg/logging"
	"github.com/mortenskoett/dotf-go/pkg/parsing"
)

// Environment used to execute commands inside
type CmdExecutor struct {
	commands map[string]Command
}

// Instantiate a new command executor to hold all available commands
func NewCmdExecutor(cmds []Command) *CmdExecutor {
	exec := CmdExecutor{
		commands: map[string]Command{},
	}
	for _, cmd := range cmds {
		exec.register(cmd)
	}
	return &exec
}

// Register a command with the executor
func (ce *CmdExecutor) register(cmd Command) error {
	_, ok := ce.commands[cmd.getName()]
	if ok {
		return &CmdAlreadyRegisteredError{cmd.getName()}
	}
	ce.commands[cmd.getName()] = cmd
	return nil
}

// A runnable command that can error at runtime
type CommandRunnable func() error

// Validate and load a command into the executor and return a runnable command or an error
func (ce *CmdExecutor) Load(
	cmdin *parsing.CommandlineInput,
	conf *parsing.DotfConfiguration,
	helpFlags []*parsing.Flag) (CommandRunnable, error) {

	if ok := userhelp(cmdin.CommandName, helpFlags); ok {
		return nil, &CmdHelpWantedError{"showing full help."}
	}

	cmd, err := parse(cmdin.CommandName, ce.commands)
	if err != nil {
		return nil, err
	}

	// Wrapped cmd.Run scoped specific to each command
	return func() error {
		// Check for command help flag
		if cmdin.Flags.OneOf(helpFlags) {
			return &CmdHelpFlagError{"help flag given", cmd}
		}

		// Check if invalid flags for current command
		var invalidflags string
		for _, cmdflag := range cmd.getFlags() {
			for _, cliflag := range cmdin.Flags.GetAllKeys() {
				if cmdflag.Name != cliflag {
					invalidflags += cliflag
					invalidflags += ", "
				}

			}
		}

		if invalidflags != "" {
			// Remove last space+comma
			invalidflags = invalidflags[:len(invalidflags)-2]
			logging.Warn("Invalid flags given for", cmd.getName(), "command:", invalidflags)
		}

		// Check for number of required positional args
		if len(cmdin.PositionalArgs) != len(cmd.getArgs()) {
			return &CmdArgumentError{fmt.Sprintf(
				"%d arguments given, but %d required.", len(cmdin.PositionalArgs), len(cmd.getArgs()))}
		}

		return cmd.Run(cmdin, conf)
	}, nil
}

// Checks whether user has inputted a request for help instead of a command name
func userhelp(cmdName string, helpFlags []*parsing.Flag) bool {
	isHelpFlagGiven := func() bool {
		for _, h := range helpFlags {
			if cmdName == "--"+h.Name { // checks w/o cmd e.g. dotf --help
				return true
			}
		}
		return false
	}

	if cmdName == "" || isHelpFlagGiven() {
		return true
	}
	return false
}

// Parses a Command name to a CommandFunc or errors
func parse(cmdName string, commands map[string]Command) (Command, error) {
	cmd, ok := commands[cmdName]
	if ok {
		return cmd, nil
	}
	return nil, &CmdUnknownCommand{fmt.Sprintf("%s command does not exist.", cmdName)}
}
