/*
Handles interaction with the command line.
*/
package terminalio

import "errors"

/* SyncLocalAndRemote will update both local repository and remote with newest changes.
If it is not possible to merge changes or if a commandline call fails, an error will be returned. */
func SyncLocalAndRemote(absPathToLocalRepo string) (bool, error) {

	// TODO
	// Return somekind of custom error for the front end to handle e.g. if a merge conflict arises

	hasNoLocalChanges, err := gitStatus.executeWithExpectedResult(absPathToLocalRepo, nothingToCommit)
	if err != nil {
		return false, err
	}

	if hasNoLocalChanges {
		return pullandMergeLatest(absPathToLocalRepo)
	}

	_, err = gitAddAll.execute(absPathToLocalRepo)
	if err != nil {
		return false, err
	}

	_, err = gitCommit.execute(absPathToLocalRepo)
	if err != nil {
		return false, err
	}

	_, err = pullandMergeLatest(absPathToLocalRepo)
	if err != nil {
		return false, err
	}

	return gitPush.executeWithExpectedResult(absPathToLocalRepo, pushWasSuccessful)
}

/* Pulls latest and attempts a merge if possible otherwise reverts the merge and returns an error. */
func pullandMergeLatest(path string) (bool, error) {
	success, err := gitPullMerge.executeWithExpectedResult(path, mergeWasSuccessful, allUpToDate)
	if err != nil {
		return false, err
	}

	if !success {
		_, err = gitAbortMerge.execute(path)
		if err != nil {
			return false, errors.New("merge was aborted, manual intervention required: " + err.Error())
		}
	}

	return true, nil
}
