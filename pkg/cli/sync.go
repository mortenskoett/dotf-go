package cli

import (
	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type syncCommand struct {
	base *CommandBase
}

func NewSyncCommand() *syncCommand {
	name := "sync"
	desc := `
	Uses local git instance to merge newest changes from git remote and then adds, commits and
	pushes latest changes to remote.`

	return &syncCommand{
		base: &CommandBase{
			Name:        name,
			Overview:    "Sync with remote using merge strategy.",
			Usage:       name + " <filepath> [--help]",
			Args:        []Arg{},
			Flags:       map[string]Arg{},
			Description: desc,
		},
	}
}

func (c *syncCommand) Base() *CommandBase {
	return c.base
}

func (c *syncCommand) Run(args *parsing.CommandLineInput, conf *parsing.DotfConfiguration) error {
	absDotfilesDir, err := terminalio.GetAndValidateAbsolutePath(conf.SyncDir)
	if err != nil {
		return err
	}

	if err := terminalio.SyncLocalRemote(absDotfilesDir); err != nil {
		return &GitError{Path: absDotfilesDir, Err: err}
	}

	return nil
}
