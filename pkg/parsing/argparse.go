package parsing

import (
	"fmt"
	"strings"
)

type CommandLineFlags struct {
	ValueFlags map[string]string // --name john or --name=john
	BoolFlags  map[string]bool   // --verbose
}

type CommandLineInput struct {
	CommandName    string   // Name of the given command
	PositionalArgs []string // command args in order
	Flags          *CommandLineFlags
}

func NewCommandLineInput() *CommandLineInput {
	return &CommandLineInput{
		CommandName:    "",
		PositionalArgs: []string{},
		Flags: &CommandLineFlags{
			ValueFlags: map[string]string{},
			BoolFlags:  map[string]bool{},
		},
	}
}

func NewCommandLineFlags() CommandLineFlags {
	return CommandLineFlags{
		ValueFlags: map[string]string{},
		BoolFlags:  map[string]bool{},
	}
}

// ParseCommandlineInput parses commands, positional arguments and flags into the CommandLineInput
// type.
func ParseCommandlineInput(osargs []string) (*CommandLineInput, error) {
	cliargs := NewCommandLineInput()

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
		return cliargs, nil
	}

	// Parse flags
	clflags, err := ParseFlags(args)
	if err != nil {
		return nil, err
	}

	cliargs.Flags = &clflags

	return cliargs, nil
}

// Parseflags parses value carrying flags and non-value carrying flags into the CommandLineFlags
// type. Regardless of the success of the function a fully functional if empty CommandlineFlags is
// returned.
func ParseFlags(args []string) (CommandLineFlags, error) {
	clflags := NewCommandLineFlags()

	if len(args) < 1 {
		return clflags, &ParseNoArgumentError{"no flags given"}
	}

	if !strings.HasPrefix(args[0], "--") {
		return clflags, &ParseNoArgumentError{"flag must be on the form --flagname"}
	}

	var value string
	for i := len(args) - 1; i >= 0; i-- {
		arg := args[i]

		// Value for current arg (flag) is set
		if value != "" {
			// Error if current flag is not a value carrier
			if !strings.HasPrefix(arg, "--") {
				return clflags, &ParseInvalidFlagError{fmt.Sprintf("given flag '%s' must be followed by a value, but was empty", arg)}
			}

			flag := strings.TrimPrefix(arg, "--")
			clflags.ValueFlags[flag] = value
			value = ""
			continue
		}

		// Arg is a boolFlag
		if strings.HasPrefix(arg, "--") {
			flag := strings.TrimPrefix(arg, "--")
			clflags.BoolFlags[flag] = true
			continue
		}

		// Arg is the value of flag parsed in the next iteration
		value = arg
	}

	return clflags, nil
}
