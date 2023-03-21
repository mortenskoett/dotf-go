package cli

import "fmt"

const (
	msgTryHelp string = "Try appending --help to see available commands."
)

type CmdHelpFlagError struct {
	message string
	Cmd     CommandPrintable
}

type CmdArgumentError struct {
	message string
}

type CmdUnknownCommand struct {
	message string
}

type CmdAlreadyRegisteredError struct {
	message string
}

type GitError struct {
	Path string
	Err  error
}

func (e *CmdUnknownCommand) Error() string {
	return fmt.Sprintf("%s %s", e.message, msgTryHelp)
}

func (e *CmdHelpFlagError) Error() string {
	return fmt.Sprint(e.message)
}

func (e *CmdArgumentError) Error() string {
	return fmt.Sprintf("%s %s", e.message, msgTryHelp)
}

func (e *CmdAlreadyRegisteredError) Error() string {
	return fmt.Sprintf("%s command already registered", e.message)
}

func (e *GitError) Error() string {
	return fmt.Sprintf("failed to execute git command in dir: %s: %v", e.Path, e.Err)
}
