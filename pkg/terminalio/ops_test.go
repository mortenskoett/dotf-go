package terminalio

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/mortenskoett/dotf-go/pkg/test"
)

// Functional tests

func TestAddFileToDotfiles(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	dfilesdir := env.DotfilesDir
	userspacedir := env.UserspaceDir

	userdirpath := "dir1/dir2"
	dir := userspacedir.AddTempDir(userdirpath)

	userspaceFile := dir.AddTempFile()

	err := AddFileToDotfiles(userspaceFile.Name(), userspacedir.Path, dfilesdir.Path)
	if err != nil {
		test.Fail(err, "No error should have happened", t)
	}

	expectedUserspaceFile := filepath.Join(dfilesdir.Path, userdirpath, filepath.Base(userspaceFile.Name()))
	expectedBackupFile := filepath.Join("/tmp/dotf-go/backups", userspaceFile.Name())

	// check if new file in dotfiles exist
	if exists, _ := checkIfFileExists(expectedUserspaceFile); !exists {
		test.Fail(exists, fmt.Sprintf("File in dotfiles dir should exist at %s", expectedUserspaceFile), t)
	}

	// check if new file at userspace location exists
	if exists, err := checkIfFileExists(userspaceFile.Name()); !exists {
		test.Fail(exists, fmt.Sprintf(
			"File in userspace dir should still exist at %s: %v", userspaceFile.Name(), err), t)
	}

	// check if new file at userspace location is symlink
	if ok, _ := isFileSymlink(userspaceFile.Name()); !ok {
		test.Fail(ok, fmt.Sprintf(
			"File in userspace dir should be a symlink at %s: %v", userspaceFile.Name(), err), t)
	}

	// check if backup exists
	if exists, err := checkIfFileExists(expectedBackupFile); !exists {
		test.Fail(exists, fmt.Sprintf(
			"File from userspace should be backed up to %s: %v", expectedBackupFile, err), t)
	}

}

