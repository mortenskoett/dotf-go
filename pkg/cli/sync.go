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
			Name:        name,
			Overview:    "Sync with remote using merge strategy.",
			Usage:       name + " <filepath> [--help]",
			Args:        []Arg{},
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

	// if args.Flags.Exists(selectFlag) {
	// 	// TODO: Implement tui selector
		// // v, _ := args.Flags.GetValue(flags.ValueFlag())
	// }

	if err := terminalio.SyncLocalRemote(absDotfilesDir); err != nil {
		return &GitError{Path: absDotfilesDir, Err: err}
	}

	return nil
}
