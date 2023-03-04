package cli

import (
	"fmt"

	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type revertCommand struct {
	commandBase
}

func NewRevertCommand(programName, commandName string) *revertCommand {
	return &revertCommand{
		commandBase{
			programName: programName,
			commandName: commandName}}
}

func (c *revertCommand) Run(args *parsing.CommandLineInput, conf *parsing.DotfConfiguration) error {
	if err := validateCliArguments(args, c); err != nil {
		return err
	}

	filepath := args.PositionalArgs[0]

	err := terminalio.RevertDotfile(filepath, conf.UserspaceDir, conf.DotfilesDir)
	if err != nil {
		return err
	}

	return nil
}

func (c *revertCommand) CmdName() string {
	return c.commandName
}

func (c *revertCommand) Overview() string {
	return "Revert file to its original location in userspace."
}

func (c *revertCommand) Arguments() []arg {
	return []arg{
		{Name: "file/dir", Description: "Path to file or dir to revert back to original location."},
	}
}

func (c *revertCommand) Usage() string {
	return fmt.Sprintf("%s %s <filepath> [--help]", c.programName, c.commandName)
}

func (c *revertCommand) Description() string {
	return `
	Will revert a file or directory previously added to dotfiles back to its original location in
	userspace. The file is moved from the dotfiles directory back to userspace where the symlink is
	removed. The command can be used both on files inside the dotfiles directory as well as symlinks
	in userspace and will do the same thing. `
}

func (c *revertCommand) ProgName() string {
	return c.programName
}
