package cli

import (
	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type addCommand struct {
	base *CommandBase
}

func NewAddCommand() *addCommand {
	name := "add"
	overview := "Move file/dir from userspace to dotfiles."
	usage := name + " <filepath> [--help]"
	args := []Arg{
		{Name: "file/dir", Description: "Path to file or dir that should be replaced by symlink."},
	}
	flags := map[string]Arg{
		"select": {
			Name:        "--select",
			Description: "Interactively select into which distros the file should be added",
		},
	}
	description := `
	Will replace a file or directory in userspace with a symlink pointing to the dotfiles directory.
	The file or the directory and its contents is copied to the dotfiles directory and a symlink is
	placed in the original location.`

	return &addCommand{
		base: &CommandBase{
			Name:        name,
			Overview:    overview,
			Description: description,
			Usage:       usage,
			Args:        args,
			Flags:       flags,
		},
	}
}

func (c *addCommand) Base() *CommandBase {
	return c.base
}

func (c *addCommand) Run(args *parsing.CommandLineInput, conf *parsing.DotfConfiguration) error {
	filepath := args.PositionalArgs[0]

	err := terminalio.AddFileToDotfiles(filepath, conf.UserspaceDir, conf.DotfilesDir)
	if err != nil {
		return err
	}

	return nil
}
