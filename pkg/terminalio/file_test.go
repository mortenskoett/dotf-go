package terminalio

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mortenskoett/dotf-go/pkg/test"
)

func TestAddFileToDotfilesNotFoundError(t *testing.T) {
	userspacefile := "adsf"
	dotfilesdir := "adsf"

	expected := &NotFoundError{}
	actual := AddFileToDotfiles(userspacefile, dotfilesdir)

	if !errors.As(actual, &expected) {
		test.Fail(actual, expected, t)
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

func TestChangeLeadingPathStrings(t *testing.T) {
	result, err := changeLeadingPath("/dotfiles/dir1/dir2/file.txt", "/dotfiles", "/userdir")
	fmt.Println(result)
	if err != nil {
		test.Fail(err, "Shouldn't fail here", t)
	}
}

func TestChangeLeadingPath(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	fromdir := env.DotfilesDir
	subfolderPath := "/bla1/bla2/"
	subfolders := fromdir.AddTempDir(subfolderPath)
	fp := subfolders.AddTempFile()

	todir := env.UserspaceDir

	result, err := changeLeadingPath(fp.Name(), fromdir.Path, todir.Path)
	if err != nil {
		test.Fail(result, err, t)
	}

	expected := filepath.Join(todir.Path, subfolderPath, filepath.Base(fp.Name()))
	if result != expected {
		test.Fail(result, expected, t)
	}
}

func TestDetachRelativePathWithStrings(t *testing.T) {
	df := "/dotfiles/d1/d2/d3/"
	bp := "/d1/d2/d3/"
	fp := bp + "file.txt"

	p, err := detachRelativePath(fp, df)

	if err != nil {
		test.Fail(err, "Shouldn't fail here", t)
	}

	if p != fp {
		test.Fail(p, bp, t)
	}
}

func TestDetachRelativePath(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	// example
	// dotfiles/d1/d2/file.txt
	// detach(dotfiles/, dotfiles/d1/d2/file.txt)
	// returns d1/d2/file.txt

	somedir := env.DotfilesDir.AddTempDir("/dotfiles/")
	basepath := somedir.AddTempDir("/bla1/bla2/")
	f := basepath.AddTempFile()

	p, err := detachRelativePath(f.Name(), basepath.Path)
	if err != nil {
		test.Fail(err, "Should not fail here", t)
	}

	// Because result has leading slash
	expected := "/" + filepath.Base(f.Name())

	// Check filename
	if p != expected {
		test.Fail(p, expected, t)
	}
}

func TestGetAndValidateAbsolutePathSame(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	f := env.UserspaceDir.AddTempFile()

	actual, err := getAbsolutePath(f.Name())

	// Check error
	if err != nil {
		test.Fail(err, "Should not fail here", t)
	}

	expected := filepath.Join(f.Name())

	// Check path -- should return the same path
	if actual != expected {
		test.Fail(actual, expected, t)
	}
}

func TestGetAndValidateAbsolutePathNotExists(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	f := "myrandomfile"

	f, err := getAndValidateAbsolutePath(f)

	if err == nil {
		test.Fail(err, "Should fail here as file does not exist.", t)
	}
}

func TestCheckIfFileExists(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	file := env.UserspaceDir.AddTempFile()

	err := checkIfFileExists(file.Name())
	if err != nil {
		test.Fail(err, "Should not fail as file exists.", t)
	}
}
