package cli

import (
	"fmt"

	"github.com/mortenskoett/dotf-go/pkg/config"
	"github.com/mortenskoett/dotf-go/pkg/logging"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

// Implements Command interface
type moveCommand struct {
	CommandBase
}

func NewMoveCommand(programName, commandName string) *moveCommand {
	return &moveCommand{
		CommandBase{
			programName: programName,
			commandName: commandName}}
}

func (c *moveCommand) Run(args *CliArguments, conf *config.DotfConfiguration) error {
	if err := checkCliArguments(args, c); err != nil {
		return err
	}

	ok := confirmByUser("\nThis operation can be desctructive. Do you want to continue?")
	if !ok {
		logging.Warn("Aborted by user")
		return nil
	}

	dotfilesDir := args.PosArgs[0]
	symlinkRootDir := args.PosArgs[1]

	err := terminalio.UpdateSymlinks(symlinkRootDir, dotfilesDir)
	if err != nil {
		return err
	}

	logging.Ok("\nAll symlinks have been updated successfully.")
	return nil
}

func (c *moveCommand) ProgName() string {
	return c.programName
}

func (c *moveCommand) CmdName() string {
	return c.commandName
}

func (c *moveCommand) Overview() string {
	return "Iterates through configs in 'dotfiles-dir' and updates matching symlinks in 'userspace-dir'."
}

func (c *moveCommand) Arguments() []Arg {
	return []Arg{
		{Name: "dotfiles-dir", Description: "Path specifies a re-located dotfiles directory."},
		{Name: "userspace-dir", Description: "Specifies userspace root directory where symlinks will be updated."},
	}
}

func (c *moveCommand) Usage() string {
	return fmt.Sprintf("%s %s <dotfiles-dir> <userspace-dir> [--help]", c.programName, c.commandName)
}

func (c *moveCommand) Description() string {
	return `
	In case the dotfiles directory has been moved, it is necessary to update all symlinks pointing
	back to the old location, to point to the new location. This application will iterate through
	all directories and files in 'from' and attempt to locate a matching symlink in the same
	location relative to the given argument 'to'. The given path 'to' is the root of the userspace,
	e.g. root of '~/' aka the home folder. Note that currently if a symlink is not found in the user
	space, then it will not be touched, however a warning will be shown.

	It is expected that the dotfiles directory has already been moved and that 'from' is the new
	location directory. 
	`
}
