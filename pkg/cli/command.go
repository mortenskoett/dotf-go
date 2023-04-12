// Package command contains handling of all dotf operations given by cli arg
package cli

import (
	"github.com/mortenskoett/dotf-go/pkg/parsing"
)

// Cli command specific flags
const (
	flagSelect string = "select"
	flagDistro string = "distro"
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
	getArgs() []arg            // Required arguments
	getFlags() []*parsing.Flag // Optional flags
	getDescription() string    // Detailed description
}

// Defines an argument for a Command
type arg struct {
	Name        string
	Description string
}

// Implements the CommandPrintable interface. Contains everything needed by a command.
type commandBase struct {
	Name        string
	Overview    string
	Usage       string
	Args        []arg
	Flags       []*parsing.Flag
	Description string
}

func (c *commandBase) getName() string {
	return c.Name
}
func (c *commandBase) getOverview() string {
	return c.Overview
}
func (c *commandBase) getUsage() string {
	return c.Usage
}
func (c *commandBase) getArgs() []arg {
	return c.Args
}
func (c *commandBase) getFlags() []*parsing.Flag {
	return c.Flags
}
func (c *commandBase) getDescription() string {
	return c.Description
}
