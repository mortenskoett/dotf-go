package terminalio

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mortenskoett/dotf-go/pkg/test"
)

/*
* Functional/integration tests implemented in terms of other terminalio function calls used for
* assertion. This requires that all units are tested individually.
 */

func Test_WriteFile(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()
	file := env.UserspaceDir.AddTempFile()
	expected := []byte("hello my friend\n")

	t.Run("File is written and backed up successfully", func(t *testing.T) {
		err := WriteFile(file.Path, expected, true)
		if err != nil {
			t.Errorf("failed running code under test: %v", err)
		}

		// Assert on file contents
		actual, err := os.ReadFile(file.Path)
		if err != nil {
			t.Errorf("failed reading actual file: %v", err)
		}

		diff := cmp.Diff(actual, expected)
		if diff != "" {
			t.Errorf("have: %+v\nwant: %+v\ndiff: %+v", actual, expected, diff)
		}

		// Assert on backup
		_, err = os.ReadFile("/tmp/dotf-go/backups/" + file.Path)
		if err != nil {
			t.Errorf("failed reading backup: %v", err)
		}
	})
}

func Test_InstallDotfile_creates_existing_symlink_in_userspace_with_overwrite(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	dfiles := env.DotfilesDir
	uspace := env.UserspaceDir

	dsomefile := dfiles.AddTempFile()
	uspaceSymlinkPath := filepath.Join(uspace.Path, filepath.Base(dsomefile.Path))

	// Create file in userspace to take up filename uspaceSymlinkPath
	fdst, err := os.Create(uspaceSymlinkPath)
	if err != nil {
		test.FailHard(err, "No error should have happened", t)
	}
	defer fdst.Close()

	// Assert file is indeed created
	exists, err := checkIfFileExists(uspaceSymlinkPath)
	if err != nil {
		test.FailHard(err, "No error should have happened", t)
	}
	if !exists {
		test.FailHardMsg("This file should at this point", exists, true, t)
	}

	err = InstallDotfile(dsomefile.Path, uspace.Path, dfiles.Path, true)
	if err != nil {
		test.FailHard(err, "No error should have happened", t)
	}

	// Assert symlink exists in uspace
	ok, err := isFileSymlink(uspaceSymlinkPath)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		test.FailHardMsg("This file should now be a symlink", ok, true, t)
	}
}

func Test_InstallDotfile_creates_existing_symlink_in_userspace_no_overwrite(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	dfiles := env.DotfilesDir
	uspace := env.UserspaceDir

	dsomefile := dfiles.AddTempFile()
	rndfile := dfiles.AddTempFile()
	uspaceSymlinkPath := filepath.Join(uspace.Path, filepath.Base(dsomefile.Path))

	// Create symlink from userspace to random useless file to take up filename uspaceSymlinkPath
	err := createSymlink(uspaceSymlinkPath, rndfile.Path)
	if err != nil {
		t.Fatal()
	}

	err = InstallDotfile(dsomefile.Path, uspace.Path, dfiles.Path, false)
	if err != nil {
		var abortErr *AbortOnOverwriteError
		if !errors.As(err, &abortErr) {
			test.FailHard(err, "No error should have happened", t)
		}
	}

	// Assert symlink exists in uspace
	ok, err := isFileSymlink(uspaceSymlinkPath)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		test.FailHardMsg("This file should be a symlink", ok, true, t)
	}
}

func Test_InstallDotfile_creates_not_existing_symlink_in_userspace_no_overwrite(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	dfiles := env.DotfilesDir
	uspace := env.UserspaceDir

	dsomefile := dfiles.AddTempFile()
	uspaceSymlinkPath := filepath.Join(uspace.Path, filepath.Base(dsomefile.Path))

	err := InstallDotfile(dsomefile.Path, uspace.Path, dfiles.Path, false)
	if err != nil {
		test.FailHard(err, "No error should have happened", t)
	}

	// Assert symlink exists in uspace
	ok, err := isFileSymlink(uspaceSymlinkPath)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		test.FailHardMsg("This file should be a symlink", ok, true, t)
	}
}

func Test_RevertDotfile_deletes_symlink_and_moves_file_back_userspace(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	dfiles := env.DotfilesDir
	uspace := env.UserspaceDir
	dsomefile := dfiles.AddTempFile()
	uspaceSymlinkPath := filepath.Join(uspace.Path, filepath.Base(dsomefile.Path))

	// Create symlink from userspace to dfiles
	err := createSymlink(uspaceSymlinkPath, dsomefile.Path)
	if err != nil {
		t.Fatal()
	}

	// Assert symlink exists
	ok, err := isFileSymlink(uspaceSymlinkPath)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		test.FailHardMsg("This file should be a symlink", ok, true, t)
	}

	// Actual test call
	err = RevertDotfile(uspaceSymlinkPath, uspace.Path, dfiles.Path)
	if err != nil {
		test.FailHard(err, "No error should have happened", t)
	}

	// Assert symlink has become a file and NOT still symlink
	ok, err = isFileSymlink(uspaceSymlinkPath)
	if err != nil {
		test.FailHard(err, "No error should have happened", t)
	}
	if ok {
		test.FailHardMsg("This file should NOT be a symlink at this point", ok, false, t)
	}

	// Assert file is not there anymore
	exists, err := checkIfFileExists(dsomefile.Path)
	if err != nil {
		test.FailHard(err, "No error should have happened", t)
	}
	if exists {
		test.FailHardMsg("This file should not exist anymore", exists, false, t)
	}
}

func Test_RevertDotfile_deletes_symlink_and_moves_file_back_dotfiles(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	dfiles := env.DotfilesDir
	uspace := env.UserspaceDir
	dsomefile := dfiles.AddTempFile()
	uspaceSymlinkPath := filepath.Join(uspace.Path, filepath.Base(dsomefile.Path))

	// Create symlink from userspace to dfiles
	err := createSymlink(uspaceSymlinkPath, dsomefile.Path)
	if err != nil {
		t.Fail()
	}

	// Assert symlink exists
	ok, err := isFileSymlink(uspaceSymlinkPath)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		test.FailHardMsg("This file should be a symlink", ok, true, t)
	}

	// Actual test call
	err = RevertDotfile(dsomefile.Path, uspace.Path, dfiles.Path)
	if err != nil {
		test.FailHard(err, "No error should have happened", t)
	}

	// Assert symlink has become a file and NOT still symlink
	ok, err = isFileSymlink(uspaceSymlinkPath)
	if err != nil {
		test.FailHard(err, "No error should have happened", t)
	}
	if ok {
		test.FailHardMsg("This file should NOT be a symlink at this point", ok, false, t)
	}

	// Assert file is not there anymore
	exists, err := checkIfFileExists(dsomefile.Path)
	if err != nil {
		test.FailHard(err, "No error should have happened", t)
	}
	if exists {
		test.FailHardMsg("This file should not exist anymore", exists, false, t)
	}
}

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
