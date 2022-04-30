package argparse

import (
	"fmt"
	"strings"

	"github.com/mortenskoett/dotf-go/pkg/cli"
	"github.com/mortenskoett/dotf-go/pkg/utils"
)

// Flags required to contain a value. These makes the parser collect: 'cmd --flag value'
type ValueFlags []string

var valueflags ValueFlags = []string{"dummy"}

// Parses the CLI input argument string. Expects complete input argument line.
func ParseCliArguments(osargs []string) (string, string, *cli.CliArguments, error) {
	execName := osargs[0]

	args := osargs[1:]

	if len(args) < 1 {
		printBasicHelp(execName)
		return "", "", nil, &ParseNoArgumentError{"no arguments given"}
	}

	cmdName := args[0]
	count := len(args)
	if cmdName == "" || cmdName == "help" || cmdName == "--help" || count == 0 {
		printFullHelp(execName)
		return "", "", nil, &ParseHelpFlagError{"showing full help."}
	}

	cmdName, cliargs, err := parse(args)
	if err != nil {
		return "", "", nil, &ParseError{fmt.Sprintf("failed to parse input: %s", err)}
	}

	return execName, cmdName, cliargs, nil
}

// Parses cli command and arguments without judgement on argument fit for Command.
func parse(osargs []string) (string, *cli.CliArguments, error) {
	cliarg := cli.NewCliArguments()

	cmdName := osargs[0]
	args := osargs[1:]

	parsePositionalInto(args, cliarg)
	parseFlagsInto(args, valueflags, cliarg)

	return cmdName, cliarg, nil
}

// Parses only positional args and stops at the first flag e.g. '--flag'. The args are added to the
// supplied cli.Arguments.
func parsePositionalInto(args []string, cliarg *cli.CliArguments) {
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			break
		} else {
			cliarg.PosArgs = append(cliarg.PosArgs, arg)
		}
	}
}

// Parses only flags but both boolean and value holding flags The flags are added to the supplied
// cli.Arguments.
func parseFlagsInto(args []string, valueflags ValueFlags, cliarg *cli.CliArguments) {
	var currentflag string

	for _, arg := range args {
		// previous arg was a value containing flag
		if currentflag != "" {
			cliarg.Flags[currentflag] = arg
			currentflag = ""
			continue
		}

		// flags
		if strings.HasPrefix(arg, "--") {
			flag := strings.ReplaceAll(arg, "--", "")

			if utils.Contains(valueflags, flag) {
				// with value
				currentflag = flag

			} else {
				// no value
				cliarg.Flags[flag] = flag
			}

		}
	}
}
