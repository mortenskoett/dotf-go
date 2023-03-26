package cli

import (
	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/parsing/flags"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type addCommand struct {
	*commandBase
}

func NewAddCommand() *addCommand {
	name := "add"
	overview := "Move file/dir from userspace to dotfiles."
	usage := name + " <filepath> [--help]"
	args := []arg{
		{name: "file/dir", description: "Path to file or dir that should be replaced by symlink."},
	}
	flags := []flag{
		{name: flags.Select, description: "Interactively select into which distros the file should be added"},
		{name: flags.Config, description: "bla bla bla"},
	}
	description := `
	Will replace a file or directory in userspace with a symlink pointing to the dotfiles directory.
	The file or the directory and its contents is copied to the dotfiles directory and a symlink is
	placed in the original location.`

	return &addCommand{
		&commandBase{
			name:        name,
			overview:    overview,
			description: description,
			usage:       usage,
			args:        args,
			flags:       flags,
		},
	}
}

func (c *addCommand) Run(args *parsing.CommandlineInput, conf *parsing.DotfConfiguration) error {
	filepath := args.PositionalArgs[0]

	if args.Flags.Exists(flags.Select) {
		// TODO: Implement tui selector
	}

	v, _ := args.Flags.GetValue(flags.Help)

	err := terminalio.AddFileToDotfiles(filepath, conf.UserspaceDir, conf.DotfilesDir)
	if err != nil {
		return err
	}

	return nil
}
