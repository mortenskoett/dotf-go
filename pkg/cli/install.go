package cli

import (
	"fmt"

	"github.com/mortenskoett/dotf-go/pkg/config"
	"github.com/mortenskoett/dotf-go/pkg/logging"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type installCommand struct {
	CommandBase
}

func NewInstallCommand(programName, commandName string) *installCommand {
	return &installCommand{
		CommandBase{
			programName: programName,
			commandName: commandName}}
}

func (c *installCommand) Run(args *CliArguments, conf *config.DotfConfiguration) error {
	if err := checkCliArguments(args, c); err != nil {
		return err
	}

	filepath := args.PosArgs[0]

	err := terminalio.InstallDotfile(filepath, conf.UserspaceDir, conf.DotfilesDir, false)
	if err != nil {
		switch e := err.(type) {
		case *terminalio.AbortOnOverwriteError:
			logging.Warn(fmt.Sprintf("A file already exist in userspace: %s", logging.Color(e.Path, logging.Green)))
			logging.Warn(fmt.Sprintf("%s needs to backup and delete this file to install the dotfile.", c.programName))

			ok := confirmByUser("Do you want to continue?")
			if ok {
				return terminalio.InstallDotfile(filepath, conf.UserspaceDir, conf.DotfilesDir, ok) // Overwrite file
			} else {
				logging.Info("Aborted by user")
				return nil
			}
		default:
			return err
		}
	}

	return nil
}

func (c *installCommand) CmdName() string {
	return c.commandName
}

func (c *installCommand) Overview() string {
	return "Install file/dir from dotfiles by creating symlink in userspace to the file."
}

func (c *installCommand) Arguments() []Arg {
	return []Arg{
		{Name: "file/dir", Description: "Path to file/dir inside dotfiles or path to file/dir in userspace."},
	}
}

func (c *installCommand) Usage() string {
	return fmt.Sprintf("%s %s <filepath> [--help]", c.programName, c.commandName)
}

func (c *installCommand) Description() string {
	return `
	Will install a file or directory from dotfiles into the same location in userspace. If an
	identically named file is found in userspace a prompt will ask whether to delete the file or
	abort the operation. The only sane thing is to remove the file in userspace as the idea is for
	the previously created dotfile to take its place.

	The command can be used both on files inside the dotfiles directory as well as files in
	userspace and will do the same thing. A file inside dotfiles can be considered the source and
	a file in userspace will be considered the target. The source must exist and the target will be
	overwritten (after a backup is made).

	- If a file from inside dotfiles is given this file will be used as installation source.
	- If a file from userspace is given this file will be used as target.
	- The command performs the same operation in both cases.
	`
}

func (c *installCommand) ProgName() string {
	return c.programName
}
