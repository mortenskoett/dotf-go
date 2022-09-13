package cli

import (
	"fmt"

	"github.com/mortenskoett/dotf-go/pkg/config"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type syncCommand struct {
	CommandBase
}

func NewSyncCommand(programName, commandName string) *syncCommand {
	return &syncCommand{
		CommandBase{
			programName: programName,
			commandName: commandName}}
}

func (c *syncCommand) Run(args *CliArguments, conf *config.DotfConfiguration) error {
	if err := checkCliArguments(args, c); err != nil {
		return err
	}

	absDotfilesDir, err := terminalio.GetAndValidateAbsolutePath(conf.DotfilesDir)
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
	return "Sync changes with remote using merge strategy if needed."
}

func (c *syncCommand) Arguments() []Arg {
	return []Arg{}
}

func (c *syncCommand) Usage() string {
	return fmt.Sprintf("%s %s <filepath> [--help]", c.programName, c.commandName)
}

func (c *syncCommand) Description() string {
	return `
	Uses local git instance to merge newest changes from git remote and then adds, commits and
	pushes latest changes to remote.
	`
}

func (c *syncCommand) ProgName() string {
	return c.programName
}
