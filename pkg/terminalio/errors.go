package terminalio

import "fmt"

/* Exported */

// The MergeFailError is returned if the remote could not be merged into the local data without interaction.
type MergeFailError struct {
	directory string
}

type UnmatchedShellReturnError struct {
	command  termCommand
	expected []commandReturn
}

type NotFoundError struct {
	path string
}

func (e *MergeFailError) Error() string {
	return fmt.Sprintf("merge was unsuccessful and rolled back (aborted). Manual intervention required in '%s'", e.directory)
}

func (e *UnmatchedShellReturnError) Error() string {
	return fmt.Sprintf("while executing '%s' in the shell, non of the following outputs where found: [%s]", e.command, e.expected)
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("file or directory was not found at: %s", e.path)
}

/* Unexported */

/* The shellExecError occurs if an unexpected error happens while executing a TermCommand in the shell. */
type shellExecError struct {
	command termCommand
}

func (e *shellExecError) Error() string {
	return fmt.Sprintf("an error has occured executing '%s' in the shell", e.command)
}
