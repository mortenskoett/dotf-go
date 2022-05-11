package cli

import (
	"fmt"

	"github.com/mortenskoett/dotf-go/pkg/config"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type revertCommand struct {
	CommandBase
}

func NewRevertCommand(programName, commandName string) *revertCommand {
	return &revertCommand{
		CommandBase{
			programName: programName,
			commandName: commandName}}
}

func (c *revertCommand) Run(args *CliArguments, conf *config.DotfConfiguration) error {
	if err := checkCliArguments(args, c); err != nil {
		return err
	}

	filepath := args.PosArgs[0]

	err := terminalio.RevertFileToUserspace(filepath, conf.HomeDir, conf.DotfilesDir)
	if err != nil {
		return err
	}

	return nil
}

func (c *revertCommand) CmdName() string {
	return c.commandName
}

func (c *revertCommand) Overview() string {
	return "reverts a file or dir from userspace to dotfiles by replacing it with a symlink and copying contents."
}

func (c *revertCommand) Arguments() []Arg {
	return []Arg{
		{Name: "file/dir", Description: "Path to file or dir that should be replaced by symlink."},
	}
}

func (c *revertCommand) Usage() string {
	return fmt.Sprintf("%s %s <filepath> [--help]", c.programName, c.commandName)
}

func (c *revertCommand) Description() string {
	return `
	Will replace a file or directory in userspace with a symlink pointing to the dotfiles directory.
	The file or the directory and its contents is copied to the dotfiles directory and a symlink is
	placed in the original location. 
	`
}

func (c *revertCommand) ProgName() string {
	return c.programName
}
