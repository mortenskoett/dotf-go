package cli

import (
	"fmt"

	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type addCommand struct {
	name string
}

func NewAddCommand() *addCommand {
	return &addCommand{name: "add"}
}

func (c *addCommand) Run(args *parsing.CommandLineInput, conf *parsing.DotfConfiguration) error {
	filepath := args.PositionalArgs[0]

	err := terminalio.AddFileToDotfiles(filepath, conf.UserspaceDir, conf.DotfilesDir)
	if err != nil {
		return err
	}

	return nil
}

func (c *addCommand) CmdName() string {
	return c.name
}

func (c *addCommand) Overview() string {
	return "Move file/dir from userspace to dotfiles."
}

func (c *addCommand) Arguments() []arg {
	return []arg{
		{name: "file/dir", description: "Path to file or dir that should be replaced by symlink."},
	}
}

func (c *addCommand) Usage() string {
	return fmt.Sprintf("%s %s <filepath> [--help]", programName, c.name)
}

func (c *addCommand) Description() string {
	return `
	Will replace a file or directory in userspace with a symlink pointing to the dotfiles directory.
	The file or the directory and its contents is copied to the dotfiles directory and a symlink is
	placed in the original location.
	`
}
