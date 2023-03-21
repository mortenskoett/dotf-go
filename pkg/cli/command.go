// Package command contains handling of all dotf operations given by cli arg
package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mortenskoett/dotf-go/pkg/logging"
	"github.com/mortenskoett/dotf-go/pkg/parsing"
)

const (
	programName string = "dotf-cli"
)

// Defines an argument for a specific Command
type arg struct {
	name        string
	description string
}

// CommandPrintable is an interface used where it is necessary to print the command details
type CommandPrintable interface {
	CmdName() string     // Name of command.
	Overview() string    // Oneliner description of the command.
	Description() string // Detailed description.
	Usage() string       // How to use the command.
	Arguments() []arg    // Required arguments in order to use the command.
}

// CommandRunner is a definition of a main operation taking a number of cli args to work on
type CommandRunner interface {
	Run(args *parsing.CommandLineInput,
		conf *parsing.DotfConfiguration) error // Run the Command using the given args and config
}

type Command interface {
	CommandPrintable
	CommandRunner
}

// Converts a slice of the runnable type to the more restrictive type. O(N).
func ConvertCommandToPrintable(cmds []Command) []CommandPrintable {
	prints := make([]CommandPrintable, len(cmds))
	for _, c := range cmds {
		prints = append(prints, c)
	}
	return prints
}

// Validates and handles the given Arguments generally against the Command and errors if not valid
// FIXME: This should be rewritten so individual commands do not depend on this
func validateCliArguments(args *parsing.CommandLineInput, c CommandPrintable) error {
	if _, ok := args.Flags.BoolFlags["help"]; ok {
		fmt.Println(generateUsage(c))
		fmt.Print("Description:")
		fmt.Println(c.Description())
		return &CmdHelpFlagError{"help flag given", c}
	}

	if len(args.PositionalArgs) != len(c.Arguments()) {
		fmt.Println(generateUsage(c))
		return &CmdArgumentError{fmt.Sprintf(
			"%d arguments given, but %d required. Try adding --help.", len(args.PositionalArgs), len(c.Arguments()))}
	}
	return nil
}

// Displays a yes/no prompt to the user and returns the boolean value of the answer
func confirmByUser(question string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		logging.Warn(question)
		logging.Input("[Y(yes)/n(no)]")

		resp, err := reader.ReadString('\n')
		if err != nil {
			logging.Fatal(err)
		}

		resp = strings.TrimSpace(resp)

		if resp == "Y" || resp == "yes" {
			return true
		} else if resp == "n" || resp == "no" {
			return false
		}
	}
}
