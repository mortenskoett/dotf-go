// Package terminalio handles interaction with the command line.
package terminalio

import (
	"log"
	"os/exec"
	"strings"
)

// termCommand is a command that can be executed in the shell.
type termCommand string

// commandReturn is the expected output to STDERR or STDOUT from executing a termCommand.
type commandReturn string

// Executes 'command' at 'path' and expects the result to contain one or more specific substrings
// 'expected'. Returns a bool and an optional error. The bool depicts whether the result contains
// any of the expected commandReturns. If the error is not nil then the boolean should be ignored.
func executeWithResult(path string, command termCommand, expected ...commandReturn) (bool, error) {
	result, err := execute(path, command)
	if err != nil {
		return false, err
	}

	for _, str := range expected {
		if strings.Contains(result, string(str)) {
			return true, nil
		}
	}
	return false, nil
}

// Executes the termCommand in the given location 'path'.
// Returns the output of the operation or an error.
// WARNING! because the command is executed as a string in the shell in order to handle
// more advaned arguments used in the called commands, this function can be used for
// malicious operations.
func execute(path string, command termCommand) (string, error) {
	args := append([]string{"-c"}, string(command)) // Prepend sh -c to give cmd as string directly to shell.
	execCmd := exec.Command("sh", args...)
	execCmd.Dir = path
	output, err := execCmd.CombinedOutput()

	log.Println("debug: ", strings.ReplaceAll(string(output), "\n", " "))

	if err != nil {
		return "", &shellExecError{command}
	}
	return string(output), nil
}
