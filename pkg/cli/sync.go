package cli

import (
	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type syncCommand struct {
	*CommandBase
}

func NewSyncCommand() *syncCommand {
	name := "sync"
	desc := `
	Uses local git instance to merge newest changes from git remote and then adds, commits and
	pushes latest changes to remote.`

	return &syncCommand{
		&CommandBase{
			name:        name,
			overview:    "Sync with remote using merge strategy.",
			usage:       name + " <filepath> [--help]",
			args:        []arg{},
			flags:       map[string]flag{},
			description: desc,
		},
	}
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
