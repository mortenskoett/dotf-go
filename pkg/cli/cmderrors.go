package cli

import "fmt"

type CmdHelpFlagError struct {
	message string
}

type CmdArgumentError struct {
	message string
}

type CmdUnknownCommand struct {
	message string
}

type GitError struct {
	Path string
	Err  error
}

func (e *CmdUnknownCommand) Error() string {
	return fmt.Sprintf(e.message)
}

func (e *CmdHelpFlagError) Error() string {
	return fmt.Sprintf(e.message)
}

func (e *CmdArgumentError) Error() string {
	return fmt.Sprintf(e.message)
}

func (e *GitError) Error() string {
	return fmt.Sprintf("failed to execute git command in dir: %s: %v", e.Path, e.Err)
}
