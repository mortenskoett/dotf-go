package cli

import (
	"github.com/mortenskoett/dotf-go/pkg/parsing"
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
		{Name: "file/dir", Description: "Path to file or dir that should be replaced by symlink."},
	}
	flags := []*parsing.Flag{
		parsing.NewFlag(flagSelect, "Interactively select individual distros into which the file should be added"),
	}
	description := `
	Will replace a file or directory in userspace with a symlink pointing to the dotfiles directory.
	The file or the directory and its contents is copied to the dotfiles directory and a symlink is
	placed in the original location.`

	return &addCommand{
		commandBase: &commandBase{
			Name:        name,
			Overview:    overview,
			Usage:       usage,
			Args:        args,
			Flags:       flags,
			Description: description,
		},
	}
}

func (c *addCommand) Run(args *parsing.CommandlineInput, conf *parsing.DotfConfiguration) error {
	filepath := args.PositionalArgs[0]

	for _, f := range c.Flags {
		switch f.Name {
		case flagSelect:
			if args.Flags.Exists(f) {
				// TODO: Implement tui selector for selectflag
			}
		}
	}

	err := terminalio.AddFileToDotfiles(filepath, conf.UserspaceDir, conf.DotfilesDir)
	if err != nil {
		return err
	}

	return nil
}
