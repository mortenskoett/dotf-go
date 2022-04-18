package cli

import "fmt"

type CmdHelpFlagError struct {
	message string
}

func (e *CmdHelpFlagError) Error() string {
	return fmt.Sprintf(e.message)
}

type CmdArgumentError struct {
	message string
}

func (e *CmdArgumentError) Error() string {
	return fmt.Sprintf(e.message)
}
