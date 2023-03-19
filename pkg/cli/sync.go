package cli

import (
	"fmt"

	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type syncCommand struct {
	commandName string
}

func NewSyncCommand(commandName string) *syncCommand {
	return &syncCommand{
		commandName: commandName,
	}
}

func (c *syncCommand) Run(args *parsing.CommandLineInput, conf *parsing.DotfConfiguration) error {
	if err := validateCliArguments(args, c); err != nil {
		return err
	}

	absDotfilesDir, err := terminalio.GetAndValidateAbsolutePath(conf.SyncDir)
	if err != nil {
		return err
	}

	if err := terminalio.SyncLocalRemote(absDotfilesDir); err != nil {
		return &GitError{Path: absDotfilesDir, Err: err}
	}

	return nil
}

func (c *syncCommand) CmdName() string {
	return c.commandName
}

func (c *syncCommand) Overview() string {
	return "Sync with remote using merge strategy."
}

func (c *syncCommand) Arguments() []arg {
	return []arg{}
}

func (c *syncCommand) Usage() string {
	return fmt.Sprintf("%s %s <filepath> [--help]", programName, c.commandName)
}

func (c *syncCommand) Description() string {
	return `
	Uses local git instance to merge newest changes from git remote and then adds, commits and
	pushes latest changes to remote.
	`
}
