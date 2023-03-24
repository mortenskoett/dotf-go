// Package command contains handling of all dotf operations given by cli arg
package cli

import (
	"github.com/mortenskoett/dotf-go/pkg/parsing"
)

const (
	programName string = "dotf-cli"
)

// Command is the dotf type denoting a runnable and printable command
type Command interface {
	CommandPrintable
	CommandRunner
}

// CommandPrintable is used where the command base info is only needed
type CommandPrintable interface {
	Base() *CommandBase // Get command base Info
}

// CommandRunner is a definition of a main operation taking a number of cli args
type CommandRunner interface {
	// Run the Command using the given args and config
	Run(args *parsing.CommandLineInput, conf *parsing.DotfConfiguration) error
}

// Defines an argument for a Command
type arg struct {
	name        string // e.g. <filepath>
	description string
}

// Defines a flag for a Command
type flag struct {
	name        string // e.g. --verbose
	description string
}

// Contains everything needed by a command. All commands must embed this struct.
type CommandBase struct {
	name        string          // Name of command.
	overview    string          // One-liner description of the command.
	usage       string          // How to use the command.
	args        []arg           // Required arguments in order to use the command.
	flags       map[string]flag // Optional command flags
	description string          // Detailed description.
}
