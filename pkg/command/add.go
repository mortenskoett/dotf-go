package command

import (
	"fmt"

	"github.com/mortenskoett/dotf-go/pkg/logger"
)

type addCommand struct {
	CommandBase
}

func NewAddCommand(programName, commandName string) *addCommand {
	return &addCommand{
		CommandBase{
			programName: programName,
			commandName: commandName}}
}

func (c *addCommand) Run(args *CliArguments) error {
	if err := checkCmdArguments(args, c); err != nil {
		return err
	}

	ok := confirmByUser("\nThis operation can be desctructive. Do you want to continue?")
	if !ok {
		logger.LogWarn("Aborted by user")
		return nil
	}

	filepath := args.PosArgs[0]
	logger.LogWarn(filepath)

	// TODO: Actual operation of command here
	// read cli arg and check that file exists in userspace
	// if not exit
	// check if file exists in dotfiles
	// if yes then write this to output and return
	// if no then
	// make backup of userspace files
	// create path in dotfiles dir
	// copy files to dotfiles dir
	// remove files from userspace
	// create symlink in userspace

	return nil
}

func (c *addCommand) CmdName() string {
	return c.commandName
}

func (c *addCommand) Overview() string {
	return "Adds a file or dir from userspace to dotfiles by replacing it with a symlink and copying contents."
}

func (c *addCommand) Arguments() *[]Arg {
	return &[]Arg{
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
