package terminalio

import (
	"log"
	"os/exec"
	"strings"
)

type termCommand struct {
	command string
}

type commandReturn string

/* Git commands. */
var (
	gitStatus     = termCommand{command: "git status"}
	gitAddAll     = termCommand{command: "git add ."}
	gitCommit     = termCommand{command: "git commit -am 'Commit made from dotf-go'"}
	gitPullMerge  = termCommand{command: "git merge origin/master -m 'Merge made by dotf-go'"}
	gitAbortMerge = termCommand{command: "git merge --abort"}
	gitPush       = termCommand{command: "git push origin master"}
)

/* Status returned by git.
Case sensitive substring contained in the returns from running commands with git version 2.32.0 */
const (
	allUpToDate             commandReturn = "Already up to date."
	nothingToCommit         commandReturn = "nothing to commit, working tree clean"
	changesToCommit         commandReturn = "Changes to be committed"
	untrackedFiles          commandReturn = "Untracked files"
	localBranchBehindRemote commandReturn = "Your branch is behind"
	canBeFastForwarded      commandReturn = "and can be fast-forwarded"
	mergeWasSuccessful      commandReturn = "Merge made by"
	pushWasSuccessful       commandReturn = "master -> master" // Something is making git push return only last line.
	branchAheadOfRemote     commandReturn = "Your branch is ahead of"
)

/* Executes the termCommand at 'path' and expects the result to contain one or more specific substrings.
Returns a bool depicting whether the result contained any of the expected substrings 'expected'.*/
func (tc *termCommand) executeExpectedResult(path string, expected ...commandReturn) (bool, error) {
	result, err := tc.execute(path)
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

/* Executes the termCommand in the given location 'path'.
Returns the output of the operation or an error.
WARNING! because the command is executed as a string in the shell to be able to handle
more advaned arguments for the called commands, this function can
be used for malicious operations. */
func (tc *termCommand) execute(path string) (string, error) {
	args := append([]string{"-c"}, tc.command)
	execCmd := exec.Command("sh", args...)
	execCmd.Dir = path
	output, err := execCmd.CombinedOutput()

	if debug_print {
		log.Println("debug: ", strings.ReplaceAll(string(output), "\n", " "))
	}

	if err != nil {
		return "", &ShellExecError{tc.command}
	}
	return string(output), nil
}
