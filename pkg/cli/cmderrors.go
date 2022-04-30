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

func (e *CmdUnknownCommand) Error() string {
	return fmt.Sprintf(e.message)
}

func (e *CmdHelpFlagError) Error() string {
	return fmt.Sprintf(e.message)
}

func (e *CmdArgumentError) Error() string {
	return fmt.Sprintf(e.message)
}
