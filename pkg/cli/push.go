package cli

import (
	"fmt"

	"github.com/mortenskoett/dotf-go/pkg/config"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type pushCommand struct {
	CommandBase
}

func NewPushCommand(programName, commandName string) *pushCommand {
	return &pushCommand{
		CommandBase{
			programName: programName,
			commandName: commandName}}
}

func (c *pushCommand) Run(args *CliArguments, conf *config.DotfConfiguration) error {
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

func (c *pushCommand) CmdName() string {
	return c.commandName
}

func (c *pushCommand) Overview() string {
	return "Merges latest changes from remote and then pushes local changes."
}

func (c *pushCommand) Arguments() []Arg {
	return []Arg{}
}

func (c *pushCommand) Usage() string {
	return fmt.Sprintf("%s %s <filepath> [--help]", c.programName, c.commandName)
}

func (c *pushCommand) Description() string {
	return `
	Uses local git instance to merge newest changes from git remote and then adds, commits and
	pushes latest changes to remote.
	`
}

func (c *pushCommand) ProgName() string {
	return c.programName
}
