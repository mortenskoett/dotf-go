// Package command contains handling of all dotf operations given by cli arg
package cli

import (
	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/parsing/flags"
)

// Program name used for printing
const (
	programName string = "dotf-cli"
)

// Command is the dotf type denoting a runnable and printable command
type Command interface {
	CommandPrintable
	CommandRunner
}

// CommandRunner is a interface to commands that can be run
type CommandRunner interface {
	Run(args *parsing.CommandlineInput, conf *parsing.DotfConfiguration) error // Run the Command using the given args and config
}

// CommandPrintable is used where the command base info is only needed
type CommandPrintable interface {
	getName() string        // Name of command
	getOverview() string    // One-liner description of the command
	getUsage() string       // How to use the command
	getArgs() []arg         // Required arguments
	getFlags() []flag       // Optional flags
	getDescription() string // Detailed description
}

// Defines an argument for a Command
type arg struct {
	name        string
	description string
}

// Defines a flag for a Command
type flag struct {
	name        flags.Flag
	description string
}

// Implements the CommandPrintable interface. Contains everything needed by a command.
type commandBase struct {
	name        string
	overview    string
	usage       string
	args        []arg
	flags       []flag
	description string
}

func (c *commandBase) getName() string {
	return c.name
}
func (c *commandBase) getOverview() string {
	return c.overview
}
func (c *commandBase) getUsage() string {
	return c.usage
}
func (c *commandBase) getArgs() []arg {
	return c.args
}
func (c *commandBase) getFlags() []flag {
	return c.flags
}
func (c *commandBase) getDescription() string {
	return c.description
}
