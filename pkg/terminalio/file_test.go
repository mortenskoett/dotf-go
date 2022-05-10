package terminalio

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mortenskoett/dotf-go/pkg/test"
)

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
	if exists, _ := CheckIfFileExists(expectedUserspaceFile); !exists {
		test.Fail(exists, fmt.Sprintf("File in dotfiles dir should exist at %s", expectedUserspaceFile), t)
	}

	// check if new file at userspace location exists
	if exists, err := CheckIfFileExists(userspaceFile.Name()); !exists {
		test.Fail(exists, fmt.Sprintf(
			"File in userspace dir should still exist at %s: %v", userspaceFile.Name(), err), t)
	}

	// check if new file at userspace location is symlink
	if ok := isFileSymlink(userspaceFile.Name()); !ok {
		test.Fail(ok, fmt.Sprintf(
			"File in userspace dir should be a symlink at %s: %v", userspaceFile.Name(), err), t)
	}

	// check if backup exists
	if exists, err := CheckIfFileExists(expectedBackupFile); !exists {
		test.Fail(exists, fmt.Sprintf(
			"File from userspace should be backed up to %s: %v", expectedBackupFile, err), t)
	}

}

func TestAddFileToDotfilesNotFoundError(t *testing.T) {
	file := "asdf"
	userspacefile := "adsf"
	dotfilesdir := "adsf"

	expected := &NotFoundError{}
	actual := AddFileToDotfiles(file, userspacefile, dotfilesdir)

	if !errors.As(actual, &expected) {
		test.Fail(actual, expected, t)
	}
}

func TestBackupFile(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	fileToBackup := env.UserspaceDir.AddTempFile().Name()
	expectedBackupPath := "/tmp/dotf-go/backups" + fileToBackup

	actual, err := BackupFile(fileToBackup)
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
	file := "/dotfiles/dir1/dir2/file.txt"
	from := "/userdir"
	to := "/dotfiles"

	result, err := ChangeLeadingPath(file, from, to)
	if err != nil {
		test.Fail(err, "Shouldn't fail here", t)
	}

	expected := filepath.Join(to, file)
	if result != expected {
		test.Fail(result, expected, t)
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

	result, err := ChangeLeadingPath(fp.Name(), fromdir.Path, todir.Path)
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

	actual, err := GetAbsolutePath(f.Name())

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

	f, err := GetAndValidateAbsolutePath(f)
	if err == nil {
		test.Fail(err, "Should fail here as file does not exist.", t)
	}
}

func TestCheckIfFileExists(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	file := env.UserspaceDir.AddTempFile()

	if exists, _ := CheckIfFileExists(file.Name()); !exists {
		test.Fail(exists, "Should not fail as file exists.", t)
	}
}

func TestCopyDir(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	/* Dir structure example:
	/tmp/userspace-dir1205530807
	/tmp/userspace-dir1205530807/inside
	/tmp/userspace-dir1205530807/inside/file-1561196035
	/tmp/userspace-dir1205530807/inside/here
	/tmp/userspace-dir1205530807/inside/here/file-4056809082
	/tmp/userspace-dir1205530807/inside/here/is
	/tmp/userspace-dir1205530807/inside/here/is/nice
	*/
	src := env.UserspaceDir
	inside := src.AddTempDir("inside")
	insideFile := inside.AddTempFile()
	here := inside.AddTempDir("here")
	hereFile := here.AddTempFile()
	here.AddTempDir("is/nice")

	dst := env.BackupDir

	res, err := copyDir(src.Path, dst.Path)
	if err != nil {
		test.Fail(err, "Should not fail here", t)
	}

	dstInsideFile, err := ChangeLeadingPath(insideFile.Name(), src.Path, dst.Path)
	if err != nil {
		test.Fail(err, "Should not fail here", t)
	}

	dstHereFile, err := ChangeLeadingPath(hereFile.Name(), src.Path, dst.Path)
	fmt.Println(dstInsideFile)
	fmt.Println(dstHereFile)
	if err != nil {
	}

	// Dst folder and content should exist
	if exists, _ := CheckIfFileExists(res); !exists {
		test.Fail(res, "File does not exist", t)
	}

	if exists, _ := CheckIfFileExists(dstInsideFile); !exists {
		test.Fail(dstInsideFile, "File does not exist", t)
	}
	if exists, _ := CheckIfFileExists(dstHereFile); !exists {
		test.Fail(dstHereFile, "File does not exist", t)
	}

	// Src folder should have been deleted
	if exists, _ := CheckIfFileExists(src.Path); exists {
		test.Fail(src.Path, "Dir should NOT exist", t)
	}
}
