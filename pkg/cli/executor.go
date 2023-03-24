package cli

import (
	"fmt"

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

// A runnable command that can error at runtime
type CommandRunnable func() error

// Validate and load a command into the executor and return a runnable command or an error
func (ce *CmdExecutor) Load(
	cliargs *parsing.CommandLineInput, config *parsing.DotfConfiguration) (CommandRunnable, error) {
	cmd, err := parse(cliargs.CommandName, ce.commands)
	if err != nil {
		return nil, err
	}

	// Wrapped command to be called at call site
	return func() error {
		err := validate(cmd, cliargs, config)
		if err != nil {
			return err
		}
		return cmd.Run(cliargs, config)
	}, nil
}

// Parses a Command name to a CommandFunc or errors
func parse(cmdName string, commands map[string]Command) (Command, error) {
	cmd, ok := commands[cmdName]
	if ok {
		return cmd, nil
	}
	return nil, &CmdUnknownCommand{fmt.Sprintf("%s command does not exist.", cmdName)}
}

// Validates the command preemptively against the given cliargs and config
func validate(c Command, args *parsing.CommandLineInput, conf *parsing.DotfConfiguration) error {
	if _, ok := args.Flags.BoolFlags["help"]; ok {
		return &CmdHelpFlagError{"help flag given", c}
	}

	if len(args.PositionalArgs) != len(c.Base().args) {
		return &CmdArgumentError{fmt.Sprintf(
			"%d arguments given, but %d required.", len(args.PositionalArgs), len(c.Base().args))}
	}
	return nil
}

// Register a command with the executor
func (ce *CmdExecutor) register(cmd Command) error {
	name := cmd.Base().name
	_, ok := ce.commands[name]
	if ok {
		return &CmdAlreadyRegisteredError{name}
	}
	ce.commands[name] = cmd
	return nil
}
