package terminalio

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mortenskoett/dotf-go/pkg/test"
)

func TestAddFileCreateSymlinkNotFoundError(t *testing.T) {
	userspacefile := "adsf"
	dotfilesdir := "adsf"

	expected := &NotFoundError{}
	actual := AddFileToDotfiles(userspacefile, dotfilesdir)

	if !errors.As(actual, &expected) {
		test.Fail(actual, expected, t)
	}
}

func TestCopyFile(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	todir := env.BackupDir
	fileToMove := env.UserspaceDir.AddTempFile().Name()
	dstFilename := "dstFileName"

	expectedPath := fmt.Sprintf("%s/%s", todir.Path, dstFilename)

	actualpath, err := copyFile(fileToMove, expectedPath)
	if err != nil {
		t.Fatal(err)
	}

	// compare returned path
	if expectedPath != actualpath {
		test.Fail(actualpath, expectedPath, t)
	}

	// check if file exists
	if _, err := os.Stat(actualpath); errors.Is(err, os.ErrNotExist) {
		test.Fail(err, expectedPath, t)
	}
}

func TestBackupFile(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	fileToBackup := env.UserspaceDir.AddTempFile().Name()
	expectedBackupPath := "/tmp/dotf-go/backups" + fileToBackup

	actual, err := backupFile(fileToBackup)
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(expectedBackupPath)

	if expectedBackupPath != actual {
		test.Fail(actual, expectedBackupPath, t)
	}
}

func TestChangeLeadingPath(t *testing.T) {
	// env := test.NewTestEnvironment()
	// defer env.Cleanup()

	// fromdir := env.DotfilesDir.AddTempDir("/bla1/bla2/")
	// f := fromdir.AddTempFile()
	// todir := env.UserspaceDir

	// actual, err := changeLeadingPath(filepath, fromdir, todir)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// defer os.Remove(expectedBackupPath)

	// if expectedBackupPath != actual {
	// 	test.Fail(actual, expectedBackupPath, t)
	// }

	// log.Println(f.Name())
	// log.Println(fromdir)

	// absPath, err := filepath.Abs(fromdir.Path)
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	log.Println(absPath)
	// }
}

func TestDetachRelativePath(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	basepath := env.DotfilesDir.AddTempDir("/bla1/bla2/")
	f := basepath.AddTempFile()

	p, err := detachRelativePath(f.Name(), basepath.Path)
	if err != nil {
		test.Fail(err, "Should not fail here", t)
	}

	// Because result has leading slash
	expected := filepath.Join("/", filepath.Base(p))

	// Check filename
	if p != expected {
		test.Fail(p, expected, t)
	}
}

func TestGetCheckAbsolutePathSame(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	f := env.UserspaceDir.AddTempFile()

	actual, err := getCheckAbsolutePath(f.Name())
	expected := filepath.Join(f.Name())

	// Check error
	if err != nil {
		test.Fail(err, "Should not fail here", t)
	}

	// Check path -- should return the same path
	if actual != expected {
		test.Fail(actual, expected, t)
	}
}

func TestGetCheckAbsolutePathDifferent(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	f := "myrandomfile"

	_, err := getCheckAbsolutePath(f)

	if err == nil {
		test.Fail(err, "Should fail here as file not found", t)
	}
}

func TestCheckIfFileExists(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	file := env.UserspaceDir.AddTempFile()

	err := checkIfFileExists(file.Name())
	if err != nil {
		test.Fail(err, "Should exist here", t)
	}
}
