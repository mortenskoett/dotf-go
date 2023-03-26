package parsing

import (
	"fmt"
	"strings"

	"github.com/mortenskoett/dotf-go/pkg/parsing/flags"
)

type CommandlineInput struct {
	CommandName    string   // Name of the given command
	PositionalArgs []string // Command args in order
	Flags          *flags.FlagHolder // Flags parsed from commandline
}

func newCommandlineInput() *CommandlineInput {
	return &CommandlineInput{
		CommandName:    "",
		PositionalArgs: []string{},
		Flags:          nil,
	}
}

// ParseCommandlineInput parses commands, positional arguments and flags into the CommandLineInput
// type.
func ParseCommandlineInput(osargs []string) (*CommandlineInput, error) {
	cliargs := newCommandlineInput()

	// Remove exec name
	args := osargs[1:]

	// Parse command
	if len(args) < 1 {
		return nil, &ParseNoArgumentError{"no command arg given"}
	}
	cliargs.CommandName = args[0]

	// Handle dotf help flag
	if cliargs.CommandName == "" ||
		cliargs.CommandName == "-h" ||
		cliargs.CommandName == "--h" ||
		cliargs.CommandName == "help" ||
		cliargs.CommandName == "--help" {
		return nil, &ParseHelpFlagError{"showing full help."}
	}

	// Remove command
	args = args[1:]

	// Parse positional args
	for _, p := range args {
		if strings.HasPrefix(p, "--") {
			// We have reached flags so no more positional args
			break
		}
		cliargs.PositionalArgs = append(cliargs.PositionalArgs, p)
	}
	posArgsCount := len(cliargs.PositionalArgs)
	args = args[posArgsCount:]

	if len(args) < 1 {
		// No flags to parse
		cliargs.Flags = flags.NewEmptyFlagHolder()
		return cliargs, nil
	}

	var err error
	cliargs.Flags, err = ParseCommandlineFlags(args)
	if err != nil {
		return nil, err
	}
	return cliargs, nil
}

func ParseCommandlineFlags(args []string) (*flags.FlagHolder, error) {
	bfs, vfs, err := parseFlags(args)
	if err != nil {
		return nil, err
	}
	return flags.NewFlagHolder(bfs, vfs), nil
}

// Parseflags parses value carrying flags and non-value carrying flags into a map of each type.
// Non-nil maps are returned regardless of errors are encountered during parsing. The given slice of
// args is iterated over in a derementing fashion in order to handle values before flagnames. This
// makes it possible to assume less about flags at this step.
func parseFlags(args []string) (bflags map[string]bool, vflags map[string]string, err error) {
	bflags = map[string]bool{}
	vflags = map[string]string{}
	err = nil

	if len(args) < 1 {
		err = &ParseNoArgumentError{"no flags given"}
		return
	}

	if !strings.HasPrefix(args[0], "--") {
		err = &ParseNoArgumentError{"flag must be on the form --flagname"}
		return
	}

	var value string
	for i := len(args) - 1; i >= 0; i-- { // The decr loop makes sure we handle values before flags
		arg := args[i]

		// Value for current arg (flag) is set
		if value != "" {
			// Error if current flag is not a value carrier
			if !strings.HasPrefix(arg, "--") {
				err = &ParseInvalidFlagError{fmt.Sprintf("given flag '%s' must be followed by a value, but was empty", arg)}
				return
			}

			flag := strings.TrimPrefix(arg, "--")
			vflags[flag] = value
			value = ""
			continue
		}

		// Arg is a boolFlag
		if strings.HasPrefix(arg, "--") {
			flag := strings.TrimPrefix(arg, "--")
			bflags[flag] = true
			continue
		}
		// Arg is the value of flag parsed in the next iteration
		value = arg
	}
	return
}
