package terminalio

type GitCommandType string

/* Git commands. */
const (
	gitStatus GitCommandType = "git status"
)

type GitReturnType string

/* Casesensitive substring contained in the returns from running commands with git version 2.32.0 */
const (
	nothingToCommit         GitReturnType = "nothing to commit"
	changesToCommit         GitReturnType = "Changes to be committed"
	untrackedFiles          GitReturnType = "Untracked files"
	localBranchBehindRemote GitReturnType = "Your branch is behind"
	canBeFastForwarded      GitReturnType = "and can be fast-forwarded"
)
