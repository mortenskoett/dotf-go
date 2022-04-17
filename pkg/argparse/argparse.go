package argparse

import (
	"fmt"
	"log"
	"strings"

	"github.com/mortenskoett/dotf-go/pkg/command"
	"github.com/mortenskoett/dotf-go/pkg/utils"
)

const (
	logo = `    _       _     __         __ _      
 __| | ___ | |_  / _|  ___  / _' | ___ 
/ _' |/ _ \|  _||  _| |___| \__. |/ _ \
\__/_|\___/ \__||_|         |___/ \___/
`
	programName string = "dotf-go"
)

var (
	// commands contains the CLI commands that are currently implemented in dotf.
	commands = map[string]command.Command{
		"add":  command.NewAddCommand(programName, "add"),
		"move": command.NewMoveCommand(programName, "move"),
	}
)

type CliArguments struct {
	command string
	posArgs []string // In order by input
	flags   map[string]string
}

// Flags required to contain a value
type ValueFlags []string

var valueflags ValueFlags = []string{"config"}

func newCliArguments() *CliArguments {
	return &CliArguments{
		flags: make(map[string]string),
	}
}

// Parses the CLI input argument string.
func HandleArguments(osargs []string) {
	args := osargs[1:]

	if len(args) < 1 {
		printBasicHelp()
		return
	}

	ops := args[0]
	count := len(args)
	if ops == "" || ops == "help" || ops == "--help" || count == 0 {
		printFullHelp()
		return
	}

	cmd, err := parseCommandName(ops)
	if err != nil {
		log.Fatal(err)
	}

	cmdargs := args[1:]
	err = cmd.Run(cmdargs)
	if err != nil {
		log.Fatal(err)
	}
}

// Parses cli arguments without judgement on fit for Command
func Parse(osargs []string) *CliArguments {
	cliarg := newCliArguments()
	args := osargs[1:]

	parsePositional(args, cliarg)
	parseFlags(args, valueflags, cliarg)

	return cliarg
}

func parseCommandName(input string) (command.Command, error) {
	cmd, ok := commands[input]
	if ok {
		return cmd, nil
	}

	return nil, fmt.Errorf("%s command does not exist. Try adding --help.", input)
}

// Parses only positional args before the first flag
func parsePositional(args []string, cliarg *CliArguments) {
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			break
		} else {
			cliarg.posArgs = append(cliarg.posArgs, arg)
		}
	}
}

// Parses only flags but both boolean and value holding flags
func parseFlags(args []string, valueflags ValueFlags, cliarg *CliArguments) {
	var currentflag string

	for _, arg := range args {

		// previous arg was a value containing flag
		if currentflag != "" {
			cliarg.flags[currentflag] = arg
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
				cliarg.flags[flag] = flag
			}

		}
	}
}
