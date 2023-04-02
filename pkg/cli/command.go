// Package command contains handling of all dotf operations given by cli arg
package cli

import (
	"github.com/mortenskoett/dotf-go/pkg/parsing"
)

// Cli command flag
const (
	flagSelect string = "select"
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
	getName() string           // Name of command
	getOverview() string       // One-liner description of the command
	getUsage() string          // How to use the command
	getArgs() []Arg            // Required arguments
	getFlags() []*parsing.Flag // Optional flags
	getDescription() string    // Detailed description
}

// Defines an argument for a Command
type Arg struct {
	Name        string
	Description string
}

// Implements the CommandPrintable interface. Contains everything needed by a command.
type CommandBase struct {
	Name        string
	Overview    string
	Usage       string
	Args        []Arg
	Flags       []*parsing.Flag
	Description string
}

func (c *CommandBase) getName() string {
	return c.Name
}
func (c *CommandBase) getOverview() string {
	return c.Overview
}
func (c *CommandBase) getUsage() string {
	return c.Usage
}
func (c *CommandBase) getArgs() []Arg {
	return c.Args
}
func (c *CommandBase) getFlags() []*parsing.Flag {
	return c.Flags
}
func (c *CommandBase) getDescription() string {
	return c.Description
}
