package terminalio

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/mortenskoett/dotf-go/pkg/test"
)

func TestAddFileCreateSymlinkNotFoundError(t *testing.T) {
	userspacefile := "adsf"
	dotfilesdir := "adsf"

	expected := &NotFoundError{}
	actual := AddFileCreateSymlink(userspacefile, dotfilesdir)

	if !errors.As(actual, &expected) {
		test.Fail(actual, expected, t)
	}
}

func TestCopyFile(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	fromdir := env.UserspaceDir
	todir := env.BackupDir
	fileToMove := fromdir.AddTempFile().Name()
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
