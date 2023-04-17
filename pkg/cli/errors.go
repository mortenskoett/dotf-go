package cli

import "fmt"

const (
	msgTryHelp string = "Try appending --help to see available commands."
)

type ErrCmdHelpWanted struct {
	message string
}

type ErrCmdHelpFlag struct {
	message string
	Cmd     CommandPrintable
}

type ErrCmdArgument struct {
	message string
}

type ErrCmdUnknownCommand struct {
	message string
}

type ErrCmdAlreadyRegistered struct {
	message string
}

type ErrGit struct {
	Path string
	Err  error
}

func (e *ErrCmdHelpWanted) Error() string {
	return fmt.Sprintf(e.message)
}

func (e *ErrCmdUnknownCommand) Error() string {
	return fmt.Sprintf("%s %s", e.message, msgTryHelp)
}

func (e *ErrCmdHelpFlag) Error() string {
	return fmt.Sprint(e.message)
}

func (e *ErrCmdArgument) Error() string {
	return fmt.Sprintf("%s %s", e.message, msgTryHelp)
}

func (e *ErrCmdAlreadyRegistered) Error() string {
	return fmt.Sprintf("%s command already registered", e.message)
}

func (e *ErrGit) Error() string {
	return fmt.Sprintf("failed to execute git command in dir: %s: %v", e.Path, e.Err)
}
