// Package command contains handling of all dotf operations given by cli arg
package cli

import (
	"bufio"
	"os"
	"strings"

	"github.com/mortenskoett/dotf-go/pkg/logging"
	"github.com/mortenskoett/dotf-go/pkg/parsing"
)

const (
	programName string = "dotf-cli"
)

// Defines an argument for a Command
type Arg struct {
	Name        string
	Description string
}

// Defines an optional flag for a command
type Flag struct {
	Name        string
	Description string
}

// Contains everything required to construct a command. All commands must embed this struct.
type CommandBase struct {
	Name        string         // Name of command.
	Overview    string         // One-liner description of the command.
	Description string         // Detailed description.
	Usage       string         // How to use the command.
	Args        []Arg          // Required arguments in order to use the command.
	Flags       map[string]Arg // Optional command flags
}

// CommandPrintable is used where the command base info is only needed
type CommandPrintable interface {
	Base() *CommandBase // Get command base Info
}

// CommandRunner is a definition of a main operation taking a number of cli args to work on
type CommandRunner interface {
	// Run the Command using the given args and config
	Run(args *parsing.CommandLineInput, conf *parsing.DotfConfiguration) error
}

type Command interface {
	CommandPrintable
	CommandRunner
}

// Converts a slice of the runnable command type to the more restrictive type. O(N).
func ConvertCommandToPrintable(cmds []Command) []CommandPrintable {
	prints := make([]CommandPrintable, 0, len(cmds))
	for _, c := range cmds {
		prints = append(prints, c)
	}
	return prints
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
