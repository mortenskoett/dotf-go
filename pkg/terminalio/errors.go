/* Contains errors used within and externally when interacting with the terminalio package. */
package terminalio

import "fmt"

type resultsNotFoundError struct {
	cmd *termCommand
}

func (r *resultsNotFoundError) Error() string {
	return fmt.Sprintf("results not found in output from command: %s", r.cmd.command)
}

type shellExecError struct {
	cmd *termCommand
}

func (s *shellExecError) Error() string {
	return fmt.Sprintf("an error has occured executing '%s' in the shell", s.cmd.command)
}
