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
func HandleArguments(osargs []string) (cli.Command, *cli.Arguments, error) {
	args := osargs[1:] // Ignore executable name

	if len(args) < 1 {
		printBasicHelp()
		return nil, nil, &ParseNoArgumentError{"no arguments given"}
	}

	cmdName := args[0]
	count := len(args)
	if cmdName == "" || cmdName == "help" || cmdName == "--help" || count == 0 {
		printFullHelp()
		return nil, nil, &ParseHelpFlagError{"showing full help."}
	}

	cmd, cliargs, err := parse(args)
	if err != nil {
		return nil, nil, &ParseError{fmt.Sprintf("failed to parse input: %s", err)}
	}

	return cmd, cliargs, nil

}

// Parses cli command and arguments without judgement on argument fit for Command
func parse(osargs []string) (cli.Command, *cli.Arguments, error) {
	cliarg := cli.NewCliArguments()

	cmdName := osargs[0]
	cmd, err := cli.ParseCommandName(cmdName)
	if err != nil {
		return nil, nil, &ParseInvalidArgumentError{fmt.Sprintf("try --help for available commands: %s", err)}
	}

	args := osargs[1:]
	parsePositional(args, cliarg)
	parseFlags(args, valueflags, cliarg)

	return cmd, cliarg, nil
}

// Parses only positional args and stops at the first flag e.g. '--flag'
func parsePositional(args []string, cliarg *cli.Arguments) {
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			break
		} else {
			cliarg.PosArgs = append(cliarg.PosArgs, arg)
		}
	}
}

// Parses only flags but both boolean and value holding flags
func parseFlags(args []string, valueflags ValueFlags, cliarg *cli.Arguments) {
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
