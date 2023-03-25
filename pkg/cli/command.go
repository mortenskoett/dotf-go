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

// CommandRunner is a interface to commands that can be run
type CommandRunner interface {
	Run(args *parsing.CommandLineInput, conf *parsing.DotfConfiguration) error // Run the Command using the given args and config
}

// CommandPrintable is used where the command base info is only needed
type CommandPrintable interface {
	getName() string           // Name of command.
	getOverview() string       // One-liner description of the command.
	getUsage() string          // How to use the command.
	getArgs() []arg            // Required arguments in order to use the command.
	getFlags() map[string]flag // Optional command flags
	getDescription() string    // Detailed description.
}

// Defines an argument for a Command
type arg struct {
	name        string
	description string
}

// Defines a flag for a Command
type flag struct {
	name        string
	description string
}

// Implements the CommandPrintable interface. Contains everything needed by a command.
type CommandBase struct {
	name        string
	overview    string
	usage       string
	args        []arg
	flags       map[string]flag
	description string
}

func (c *CommandBase) getName() string {
	return c.name
}
func (c *CommandBase) getOverview() string {
	return c.overview
}
func (c *CommandBase) getUsage() string {
	return c.usage
}
func (c *CommandBase) getArgs() []arg {
	return c.args
}
func (c *CommandBase) getFlags() map[string]flag {
	return c.flags
}
func (c *CommandBase) getDescription() string {
	return c.description
}
