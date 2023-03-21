package cli

import (
	"fmt"
	"sort"

	"github.com/mortenskoett/dotf-go/pkg/parsing"
)

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

// Load a command into the executor and return a runnable command or an error
func (ce *CmdExecutor) Load(
	cliargs *parsing.CommandLineInput, config *parsing.DotfConfiguration) (CommandRunnable, error) {
	cmd, err := parse(cliargs.CommandName, ce.commands)
	if err != nil {
		return nil, &CmdUnknownCommand{fmt.Sprintf("try --help for available commands: %s", err)}
	}

	return func() error {
		err := validate(cmd, cliargs, config)
		if err != nil {
			return err
		}
		return cmd.Run(cliargs, config)
	}, nil
}

// Validates the command preemptively against the given cliargs and config
func validate(c Command, args *parsing.CommandLineInput, conf *parsing.DotfConfiguration) error {
	// validate config specifically against this command
	// validate flags specifically against this command
	// on error return typed error so main can print help etc
	if _, ok := args.Flags.BoolFlags["help"]; ok {
		return &CmdHelpFlagError{"help flag given"}
	}

	if len(args.PositionalArgs) != len(c.Arguments()) {
		return &CmdArgumentError{fmt.Sprintf(
			"%d arguments given, but %d required. Try adding --help.", len(args.PositionalArgs), len(c.Arguments()))}
	}

	return nil
}

// Get commands as sorted slice for pretty printing
func (ce *CmdExecutor) GetPrintableCmds() []CommandPrintable {
	cmds := make([]CommandPrintable, 0, len(ce.commands))
	for _, cmd := range ce.commands {
		cmds = append(cmds, cmd)
	}
	sort.SliceStable(cmds, func(i, j int) bool {
		return cmds[i].CmdName() < cmds[j].CmdName()
	})
	return cmds
}

// Parses a Command name to a CommandFunc or errors
func parse(cmdName string, commands map[string]Command) (Command, error) {
	cmd, ok := commands[cmdName]
	if ok {
		return cmd, nil
	}
	return nil, &CmdArgumentError{fmt.Sprintf("%s command does not exist.", cmdName)}
}

// Register a command with the executor
func (ce *CmdExecutor) register(cmd Command) error {
	_, ok := ce.commands[cmd.CmdName()]
	if ok {
		return &CmdAlreadyRegisteredError{cmd.CmdName()}
	}
	ce.commands[cmd.CmdName()] = cmd
	return nil
}
