package cli

import (
	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type revertCommand struct {
	*commandBase
}

func NewRevertCommand() *revertCommand {
	name := "revert"
	desc := `
	Will revert a file or directory previously added to dotfiles back to its original location in
	userspace. The file is moved from the dotfiles directory back to userspace where the symlink is
	removed. The command can be used both on files inside the dotfiles directory as well as symlinks
	in userspace and will do the same thing. `

	return &revertCommand{
		&commandBase{
			Name:        name,
			Overview:    "Revert file to its original location in userspace.",
			Usage:       name + " <filepath> [--help]",
			Args:        []arg{{Name: "file/dir", Description: "Path to file or dir to revert back to original location."}},
			Flags:       []*parsing.Flag{},
			Description: desc,
		},
	}
}

func (c *revertCommand) Run(args *parsing.CommandlineInput, conf *parsing.DotfConfiguration) error {
	filepath := args.PositionalArgs[0]

	err := terminalio.RevertDotfile(filepath, conf.UserspaceDir, conf.DotfilesDir)
	if err != nil {
		return err
	}

	return nil
}
