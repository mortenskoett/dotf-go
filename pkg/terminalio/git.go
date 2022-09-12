package terminalio

// Git commands.
const (
	gitStatus    termCommand = "git status"
	gitAddAll    termCommand = "git add ."
	gitCommit    termCommand = "git commit -am 'Commit made from dotf-go'"
	gitFetch     termCommand = "git fetch"
	gitPullMerge termCommand = "git merge origin/master -m 'Merge made by dotf-go'"
	// gitPullMerge  termCommand = "git pull --no-edit" // Merge is default
	gitAbortMerge termCommand = "git merge --abort"
	gitPush       termCommand = "git push origin master"
)

// Status returned by git. Case sensitive substring contained in the first line of the returns from
// running commands with git version 2.32.0
const (
	allUpToDate     commandReturn = "Already up to date."
	nothingToCommit commandReturn = "nothing to commit, working tree clean"
	mergeSuccess    commandReturn = "Merge made by"
	pushSuccess     commandReturn = "master -> master" // Something is making git push return only last line.
)

// SyncLocalRemote uses Git to update local and remote repository with newest changes from either place.
// The given path 'absPathToLocalRepo' must point to a directory initialized with git and with push/pull
// abilities to a remote.
// If it is not possible to merge changes or if a command fails in the shell, an error will be returned.
func SyncLocalRemote(absPath string) error {
	hasNoLocalChanges, err := executeWithResult(absPath, gitStatus, nothingToCommit)
	if err != nil {
		return err
	}

	_, err = execute(absPath, gitFetch)
	if err != nil {
		return err
	}

	if hasNoLocalChanges {
		return pullMerge(absPath)
	}

	err = addCommitAll(absPath)
	if err != nil {
		return err
	}

	err = pullMerge(absPath)
	if err != nil {
		return err
	}

	expected := pushSuccess
	found, err := executeWithResult(absPath, gitPush, expected)
	if err != nil {
		return err
	}

	if !found {
		return &UnmatchedShellReturnError{gitPush, []commandReturn{expected}}
	}

	return nil
}

// Simply stages everything and creates a combined commit.
func addCommitAll(path string) error {
	_, err := execute(path, gitAddAll)
	if err != nil {
		return err
	}

	_, err = execute(path, gitCommit)
	if err != nil {
		return err
	}
	return nil
}

// Pulls latest and attempts a merge if possible otherwise reverts the merge and returns an error.
func pullMerge(path string) error {
	success, err := executeWithResult(path, gitPullMerge, mergeSuccess, allUpToDate)
	if err != nil {
		return err
	}

	if !success {
		_, err = execute(path, gitAbortMerge)
		if err != nil {
			return &MergeFailError{path}
		}
	}

	return nil
}
