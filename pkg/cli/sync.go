package cli

import (
	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type syncCommand struct {
	*commandBase
}

func NewSyncCommand() *syncCommand {
	name := "sync"
	desc := `
	Uses local git instance to merge newest changes from git remote and then adds, commits and
	pushes latest changes to remote.`

	return &syncCommand{
		&commandBase{
			Name:        name,
			Overview:    "Sync with remote using merge strategy.",
			Usage:       name + " [--help]",
			Args:        []arg{},
			Flags:       []*parsing.Flag{},
			Description: desc,
		},
	}
}

func (c *syncCommand) Run(args *parsing.CommandlineInput, conf *parsing.DotfConfiguration) error {
	absDotfilesDir, err := terminalio.GetAndValidateAbsolutePath(conf.SyncDir)
	if err != nil {
		return err
	}

	if err := terminalio.SyncLocalRemote(absDotfilesDir); err != nil {
		return &ErrGit{Path: absDotfilesDir, Err: err}
	}

	return nil
}
