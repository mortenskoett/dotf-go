package terminalio

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/mortenskoett/dotf-go/pkg/test"
)

/* Functional tests implemented in terms of other terminalio function calls used for assertion. This
* requires that all units are tested individually. */

func Test_AddFileToDotfiles_unknown_path_gives_error(t *testing.T) {
	file := "asdf"
	userspacefile := "adsf"
	dotfilesdir := "adsf"

	expected := &FileNotFoundError{}
	actual := AddFileToDotfiles(file, userspacefile, dotfilesdir)

	if !errors.As(actual, &expected) {
		test.Fail(actual, expected, t)
	}
}

func Test_AddFileToDotfiles_successfully_copies_file_creates_symlink(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	dfilesdir := env.DotfilesDir
	userspacedir := env.UserspaceDir

	userdirpath := "dir1/dir2"
	dir := userspacedir.AddTempDir(userdirpath)

	userspaceFile := dir.AddTempFile()

	// Function under test
	err := AddFileToDotfiles(userspaceFile.Path, userspacedir.Path, dfilesdir.Path)
	if err != nil {
		test.Fail(err, "No error should have happened", t)
	}

	expectedUserspaceFile := filepath.Join(dfilesdir.Path, userdirpath, filepath.Base(userspaceFile.Path))
	expectedBackupFile := filepath.Join("/tmp/dotf-go/backups", userspaceFile.Path)

	// check if new file in dotfiles exist
	if exists, _ := checkIfFileExists(expectedUserspaceFile); !exists {
		test.Fail(exists, fmt.Sprintf("File in dotfiles dir should exist at %s", expectedUserspaceFile), t)
	}

	// check if new file at userspace location exists
	if exists, err := checkIfFileExists(userspaceFile.Path); !exists {
		test.Fail(exists, fmt.Sprintf(
			"File in userspace dir should still exist at %s: %v", userspaceFile.Path, err), t)
	}

	// check if new file at userspace location is symlink
	if ok, _ := isFileSymlink(userspaceFile.Path); !ok {
		test.Fail(ok, fmt.Sprintf(
			"File in userspace dir should be a symlink at %s: %v", userspaceFile.Path, err), t)
	}

	// check if backup exists
	if exists, err := checkIfFileExists(expectedBackupFile); !exists {
		test.Fail(exists, fmt.Sprintf(
			"File from userspace should be backed up to %s: %v", expectedBackupFile, err), t)
	}
}
