package terminalio

/* Git commands. */
const (
	gitStatus     termCommand = "git status"
	gitAddAll     termCommand = "git add ."
	gitCommit     termCommand = "git commit -am 'Commit made from dotf-go'"
	gitPullMerge  termCommand = "git merge origin/master -m 'Merge made by dotf-go'"
	gitAbortMerge termCommand = "git merge --abort"
	gitPush       termCommand = "git push origin master"
)

/* Status returned by git.
Case sensitive substring contained in the returns from running commands with git version 2.32.0 */
const (
	allUpToDate        commandReturn = "Already up to date."
	nothingToCommit    commandReturn = "nothing to commit, working tree clean"
	mergeWasSuccessful commandReturn = "Merge made by"
	pushWasSuccessful  commandReturn = "master -> master" // Something is making git push return only last line.
)

/* SyncLocalAndRemote will update both local repository and remote with newest changes.
If it is not possible to merge changes or if a commandline call fails, an error will be returned. */
func SyncLocalAndRemote(absPathToLocalRepo string) error {
	hasNoLocalChanges, err := executeExpectedResult(gitStatus, absPathToLocalRepo, nothingToCommit)
	if err != nil {
		return err
	}

	if hasNoLocalChanges {
		return pullMergeLatest(absPathToLocalRepo)
	}

	err = addCommitChanges(absPathToLocalRepo)
	if err != nil {
		return err
	}

	err = pullMergeLatest(absPathToLocalRepo)
	if err != nil {
		return err
	}

	expected := pushWasSuccessful
	contains, err := executeExpectedResult(gitPush, absPathToLocalRepo, expected)
	if err != nil {
		return err
	}

	if !contains {
		return &UnmatchedShellReturnError{gitPush, []commandReturn{expected}}
	}

	return nil
}

/* Simply stages everything and creates a combined commit. */
func addCommitChanges(path string) error {
	_, err := execute(gitAddAll, path)
	if err != nil {
		return err
	}

	_, err = execute(gitCommit, path)
	if err != nil {
		return err
	}
	return nil
}

/* Pulls latest and attempts a merge if possible otherwise reverts the merge and returns an error. */
func pullMergeLatest(path string) error {
	success, err := executeExpectedResult(gitPullMerge, path, mergeWasSuccessful, allUpToDate)
	if err != nil {
		return err
	}

	if !success {
		_, err = execute(gitAbortMerge, path)
		if err != nil {
			return &MergeFailError{path}
		}
	}

	return nil
}
