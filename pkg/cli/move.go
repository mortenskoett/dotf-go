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

	ok := confirmByUser("This operation can be desctructive. Do you want to continue?")
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
	return "Iterate over files inside dotfiles and update matching symlinks in userspace."
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
	In case the dotfiles directory has been moved, here denoted by given 'dotfiles-dir', it is
	necessary to update all symlinks pointing back to the previous location, to point to the new
	location.
	This command will iterate through all directories and files in given 'dotfiles-dir', and attempt
	to locate a matching symlink in the same location relative to the given argument but in given
	'userspace-dir'. The path 'userspace-dir' must be the root of the configured userspace, e.g.
	'~/' aka the home folder. Note that currently if a symlink is not found in userspace, then it
	will not be touched, however a warning will be shown.

	It is expected that the dotfiles directory has already been moved and that 'dotfiles-dir' is the
	new location directory.
	`
}
