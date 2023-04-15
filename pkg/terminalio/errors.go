package terminalio

import "fmt"

/* Exported */

// The ErrMergeFail is returned if the remote could not be merged into the local data without interaction.
type ErrMergeFail struct {
	directory string
}

type ErrUnmatchedShellReturn struct {
	command  termCommand
	expected []commandReturn
}

type ErrFileNotFound struct {
	path string
}

type ErrFileAlreadyExists struct {
	path string
}

type ErrSymlinkNotFound struct {
	path string
}

type ErrAbortOnOverwrite struct {
	Path string
}

func (e *ErrAbortOnOverwrite) Error() string {
	return fmt.Sprintf("file or directory was present at location: %s. User interaction required.", e.Path)
}

func (e *ErrFileAlreadyExists) Error() string {
	return fmt.Sprintf("file or directory was already present at location: %s", e.path)
}

func (e *ErrMergeFail) Error() string {
	return fmt.Sprintf("merge was unsuccessful and rolled back (aborted). Manual intervention required in '%s'", e.directory)
}

func (e *ErrUnmatchedShellReturn) Error() string {
	return fmt.Sprintf("while executing '%s' in the shell, non of the following outputs where found: [%s]", e.command, e.expected)
}

func (e *ErrFileNotFound) Error() string {
	return fmt.Sprintf("file or directory was not found at: %s", e.path)
}

func (e *ErrSymlinkNotFound) Error() string {
	return fmt.Sprintf("file or directory was not a symlink at: %s", e.path)
}

/* Unexported */

/* The errShellExec occurs if an unexpected error happens while executing a TermCommand in the shell. */
type errShellExec struct {
	command termCommand
	output  string
}

func (e *errShellExec) Error() string {
	return fmt.Sprintf("an error has occured executing '%s' in the shell: %s", e.command, e.output)
}
