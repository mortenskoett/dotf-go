package cli

import (
	"fmt"

	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type addCommand struct {
	commandBase
}

func NewAddCommand(programName, commandName string) *addCommand {
	return &addCommand{
		commandBase{
			programName: programName,
			commandName: commandName}}
}

func (c *addCommand) Run(args *parsing.CommandLineInput, conf *parsing.DotfConfiguration) error {
	if err := validateCliArguments(args, c); err != nil {
		return err
	}

	filepath := args.PositionalArgs[0]

	err := terminalio.AddFileToDotfiles(filepath, conf.UserspaceDir, conf.DotfilesDir)
	if err != nil {
		return err
	}

	return nil
}

func (c *addCommand) CmdName() string {
	return c.commandName
}

func (c *addCommand) Overview() string {
	return "Move file/dir from userspace to dotfiles."
}

func (c *addCommand) Arguments() []arg {
	return []arg{
		{Name: "file/dir", Description: "Path to file or dir that should be replaced by symlink."},
	}
}

func (c *addCommand) Usage() string {
	return fmt.Sprintf("%s %s <filepath> [--help]", c.programName, c.commandName)
}

func (c *addCommand) Description() string {
	return `
	Will replace a file or directory in userspace with a symlink pointing to the dotfiles directory.
	The file or the directory and its contents is copied to the dotfiles directory and a symlink is
	placed in the original location.
	`
}

func (c *addCommand) ProgName() string {
	return c.programName
}
