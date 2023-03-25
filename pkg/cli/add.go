package cli

import (
	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type addCommand struct {
	*CommandBase
}

func NewAddCommand() *addCommand {
	name := "add"
	overview := "Move file/dir from userspace to dotfiles."
	usage := name + " <filepath> [--help]"
	args := []arg{
		{name: "file/dir", description: "Path to file or dir that should be replaced by symlink."},
	}
	flags := map[string]flag{
		"select": {
			name:        "--select",
			description: "Interactively select into which distros the file should be added",
		},
	}
	description := `
	Will replace a file or directory in userspace with a symlink pointing to the dotfiles directory.
	The file or the directory and its contents is copied to the dotfiles directory and a symlink is
	placed in the original location.`

	return &addCommand{
		&CommandBase{
			name:        name,
			overview:    overview,
			description: description,
			usage:       usage,
			args:        args,
			flags:       flags,
		},
	}
}

func (c *addCommand) Run(args *parsing.CommandLineInput, conf *parsing.DotfConfiguration) error {
	filepath := args.PositionalArgs[0]

	err := terminalio.AddFileToDotfiles(filepath, conf.UserspaceDir, conf.DotfilesDir)
	if err != nil {
		return err
	}

	return nil
}
