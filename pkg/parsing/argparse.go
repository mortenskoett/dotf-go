package parsing

import (
	"fmt"
	"strings"
)

type CommandlineInput struct {
	CommandName    string      // Name of the given command
	PositionalArgs []string    // Command args in order
	Flags          *FlagHolder // Flags parsed from commandline
}

func newCommandlineInput() *CommandlineInput {
	return &CommandlineInput{
		CommandName:    "",
		PositionalArgs: []string{},
		Flags:          &FlagHolder{},
	}
}

// ParseCommandlineArgs parses commands, positional arguments and flags into the CommandLineInput
// type.
func ParseCommandlineArgs(osargs []string) (*CommandlineInput, error) {
	// Remove exec name
	args := osargs[1:]

	// Parse command
	if len(args) < 1 {
		return nil, &ParseNoArgumentError{"no command arg given"}
	}

	cliargs := newCommandlineInput()
	cliargs.CommandName = args[0]

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
		cliargs.Flags = NewEmptyFlagHolder()
		return cliargs, nil
	}

	fmap, err := parseFlags(args)
	if err != nil {
		return nil, err
	}

	cliargs.Flags = NewFlagHolder(fmap)
	return cliargs, nil
}

func ParseCommandlineFlags(args []string) (*FlagHolder, error) {
	fmap, err := parseFlags(args)
	if err != nil {
		return nil, err
	}
	return NewFlagHolder(fmap), nil
}

// Parseflags parses value carrying flags and non-value carrying flags into a single map with value
// carrying flags pointing to a string and boolean based flags pointing to an empty string.
// The given slice of args is iterated over in a decrementing fashion in order to handle values
// before flagnames. This makes it possible to assume less about flags at this step.
func parseFlags(args []string) (map[string]string, error) {
	pflags := make(map[string]string, len(args))

	if len(args) < 1 {
		return nil, &ParseNoArgumentError{"no flags given"}
	}

	if !strings.HasPrefix(args[0], "--") {
		return nil, &ParseNoArgumentError{"flag must be on the form --flagname"}
	}

	var value string
	for i := len(args) - 1; i >= 0; i-- { // The decr loop makes sure we handle values before flags
		arg := args[i]

		// Value for current arg (flag) is set
		if value != "" {
			// Error if current flag is not a value carrier
			if !strings.HasPrefix(arg, "--") {
				return nil, &ParseInvalidFlagError{fmt.Sprintf("given flag '%s' must be followed by a value, but was empty", arg)}
			}
			flag := strings.TrimPrefix(arg, "--")
			if exists(flag, pflags) {
				return nil, &ParseInvalidFlagError{fmt.Sprintf("given flag '%s' was encountered twice", flag)}
			}

			pflags[flag] = value
			value = ""
			continue
		}

		// Arg is a boolFlag
		if strings.HasPrefix(arg, "--") {
			flag := strings.TrimPrefix(arg, "--")
			if exists(flag, pflags) {
				return pflags, &ParseInvalidFlagError{fmt.Sprintf("given flag '%s' was encountered twice", flag)}
			}
			pflags[flag] = ""
			continue
		}
		// Arg is the value of flag parsed in the next iteration
		value = arg
	}
	return pflags, nil
}

func exists[T comparable, K any](v T, m map[T]K) bool {
	if _, ok := m[v]; ok {
		return ok
	}
	return false
}

