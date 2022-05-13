package cli

import (
	"fmt"

	"github.com/mortenskoett/dotf-go/pkg/config"
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

	// filepath := args.PosArgs[0]

	// TODO
	// err := terminalio.installDotfile(filepath, conf.HomeDir, conf.DotfilesDir)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (c *installCommand) CmdName() string {
	return c.commandName
}

func (c *installCommand) Overview() string {
	return "Installs a dotfile into the userspace by symlinking to it from userspace."
}

func (c *installCommand) Arguments() []Arg {
	return []Arg{
		{Name: "file/dir", Description: "Path to file or dir inside dotfiles to install or path to userspace file to overwrite."},
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
	the previously created dotfile to take its place. The command can be used both on files inside
	the dotfiles directory as well as files in userspace and will do the same thing. `
}

func (c *installCommand) ProgName() string {
	return c.programName
}
