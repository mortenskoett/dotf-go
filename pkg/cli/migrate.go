package cli

import (
	"github.com/mortenskoett/dotf-go/pkg/logging"
	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

// Implements Command interface
type migrateCommand struct {
	*commandBase
}

func NewMigrateCommand() *migrateCommand {
	name := "migrate"
	desc := `
	In case the dotfiles directory has been moved, here denoted by given 'dotfiles-dir', it is
	necessary to migrate all symlinks pointing back to the previous location, to point to the new
	location.

	This command will iterate through all directories and files in given 'dotfiles-dir', and attempt
	to locate a matching symlink in the same location relative to the given argument but in given
	'userspace-dir'. The path 'userspace-dir' must be the root of the configured userspace, e.g.
	'~/' aka the home folder. Note that currently if a symlink is not found in userspace, then it
	will not be touched, however a warning will be shown.

	It is expected that the dotfiles directory has already been moved and that 'dotfiles-dir' is the
	new location directory.`

	return &migrateCommand{
		&commandBase{
			Name:     name,
			Overview: "Migrate userspace symlinks on dotfiles dir location change.",
			Usage:    name + " <dotfiles-dir> <userspace-dir> [--help]",
			Args: []arg{
				{Name: "dotfiles-dir", Description: "Path specifies a re-located dotfiles directory."},
				{Name: "userspace-dir", Description: "Specifies userspace root directory where symlinks will be updated."},
			},
			Flags:       []*parsing.Flag{},
			Description: desc,
		},
	}
}

func (c *migrateCommand) Run(args *parsing.CommandlineInput, conf *parsing.DotfConfiguration) error {
	ok := confirmByUser("This operation can be desctructive. Do you want to continue?")
	if !ok {
		logging.Warn("Aborted by user")
		return nil
	}

	dotfilesDir := args.PositionalArgs[0]
	symlinkRootDir := args.PositionalArgs[1]

	err := terminalio.UpdateSymlinks(symlinkRootDir, dotfilesDir)
	if err != nil {
		return err
	}

	logging.Ok("\nAll symlinks have been updated successfully.")
	return nil
}

