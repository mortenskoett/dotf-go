// Package terminalio handles interaction with the command line.
package terminalio

import (
	"os/exec"
	"strings"

	"github.com/mortenskoett/dotf-go/pkg/logging"
)

// termCommand is a command that can be executed in the shell.
type termCommand string

// commandReturn is the expected output to STDERR or STDOUT from executing a termCommand.
type commandReturn string

// Executes 'command' at 'path' and expects the result to contain one or more specific substrings
// 'expected'. Returns a bool and an optional error. The bool depicts whether the result contains
// any of the expected commandReturns. If the error is not nil then the boolean should be ignored.
// Even if an error is returned the boolean will still describe whether the output from git
// contained the expected.
func executeWithResult(path string, command termCommand, expected ...commandReturn) (bool, error) {
	result, err := execute(path, command)
	var asExpected bool

	for _, str := range expected {
		if strings.Contains(result, string(str)) {
			asExpected = true
		}
	}

	if err != nil {
		return asExpected, err
	}

	return asExpected, nil
}

// Executes the termCommand in the given location 'path'.
// Returns the output of the operation or an error.
// WARNING! Because the command is executed as a string in the shell in order to handle
// more advaned arguments used in the called commands, this function can be used for
// malicious operations.
func execute(path string, command termCommand) (string, error) {
	args := append([]string{"-c"}, string(command)) // Prepend sh -c to give cmd as string directly to shell.
	execCmd := exec.Command("sh", args...)
	execCmd.Dir = path
	output, err := execCmd.CombinedOutput()

	// Show command output in terminal
	logging.Info(logging.Color(string(command), logging.Yellow))
	if len(output) != 0 {
		logging.Info(logging.Color(string(output), logging.Green))
	}

	if err != nil {
		return "", &errShellExec{command, string(output)}
	}
	return string(output), nil
}
