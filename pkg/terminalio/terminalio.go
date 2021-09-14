/*
Handles interaction with the command line.
*/
package terminalio

const (
	debug_print bool = true
)

// TODO: Remove these returned booleans. It doesn't really make a lot of sense.

/* SyncLocalAndRemote will update both local repository and remote with newest changes.
If it is not possible to merge changes or if a commandline call fails, an error will be returned. */
func SyncLocalAndRemote(absPathToLocalRepo string) (bool, error) {
	hasNoLocalChanges, err := gitStatus.executeExpectedResult(absPathToLocalRepo, nothingToCommit)
	if err != nil {
		return false, err
	}

	if hasNoLocalChanges {
		return pullMergeLatest(absPathToLocalRepo)
	}

	err = addCommitChanges(absPathToLocalRepo)
	if err != nil {
		return false, err
	}

	_, err = pullMergeLatest(absPathToLocalRepo)
	if err != nil {
		return false, err
	}

	return gitPush.executeExpectedResult(absPathToLocalRepo, pushWasSuccessful)
}

/* Simply stages everything and creates a combined commit. */
func addCommitChanges(path string) error {
	_, err := gitAddAll.execute(path)
	if err != nil {
		return err
	}

	_, err = gitCommit.execute(path)
	if err != nil {
		return err
	}
	return nil
}

/* Pulls latest and attempts a merge if possible otherwise reverts the merge and returns an error. */
func pullMergeLatest(path string) (bool, error) {
	success, err := gitPullMerge.executeExpectedResult(path, mergeWasSuccessful, allUpToDate)
	if err != nil {
		return false, err
	}

	if !success {
		_, err = gitAbortMerge.execute(path)
		if err != nil {
			return false, &MergeFailError{path}
		}
	}

	return true, nil
}
