/* Contains errors used within and externally when interacting with the terminalio package. */
package terminalio

import "fmt"

type ShellExecError struct {
	command string
}

func (s *ShellExecError) Error() string {
	return fmt.Sprintf("an error has occured executing '%s' in the shell", s.command)
}

type MergeFailError struct {
	directory string
}

func (m *MergeFailError) Error() string {
	return fmt.Sprintf("merge was unsuccessful and rolled back (aborted). Manual intervention required in '%s'", m.directory)
}
