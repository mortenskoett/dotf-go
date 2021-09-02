/*
Handles interaction with the command line.
*/
package terminalio

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

/* SyncLocalAndRemote will update both local repository and remote with newest changes.
If it is not possible to merge changes or a command line call fails, an error will be returned. */
func SyncLocalAndRemote(absPathToLocalRepo string) (bool, error) {
	status, err := executeCommand(gitStatus, absPathToLocalRepo)
	if err != nil {
		return false, err
	}

	fmt.Println("info:", status)

	// TODO
	// Return somekind of custom error for the front end to handle

	// TODO
	// 1. Commit local changes
	// 2. Fetch status from server
	// 2. 	If OK then push()
	// 2. 	If DIFF then pull merge
	// 2. 		if OK then push
	// 2.			if CONFLICT then a) rollback merge + b) advise user somehow

	return true, nil
}

func executeCommand(command GitCommandType, absDirPath string) (string, error) {
	cmd, args := SeparateCommandAndArguments(command)
	execCmd := exec.Command(cmd, args)
	execCmd.Dir = absDirPath
	output, err := execCmd.Output()
	if err != nil {
		log.Println("an error has occured running: ", cmd, args)
		return "", err
	}
	return string(output), nil
}

func SeparateCommandAndArguments(command GitCommandType) (string, string) {
	cmdAndArgs := strings.SplitN(string(command), " ", 2)
	if len(cmdAndArgs) < 2 {
		return cmdAndArgs[0], ""
	}
	return cmdAndArgs[0], cmdAndArgs[1]
}
