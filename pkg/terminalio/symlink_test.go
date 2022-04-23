package terminalio

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/mortenskoett/dotf-go/pkg/test"
)

func TestAddFileNotFoundError(t *testing.T) {
	userspacefile := "adsf"
	dotfilesdir := "adsf"

	expected := &NotFoundError{}
	actual := AddFile(userspacefile, dotfilesdir)

	if !errors.As(actual, &expected) {
		test.Fail(actual, expected, t)
	}
}

func TestBackupFileTemp(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	fromdir := env.UserspaceDir
	todir := env.BackupDir

	fileToMove := fromdir.AddTempFile().Name()

	err := backupFileTemp(fileToMove, todir.Path)

	files, err := ioutil.ReadDir(todir.Path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}

	if err != nil {
		test.Fail(err, "should not fail", t)
	}

	expectedPath := fmt.Sprintf("%s%s", todir.Name, fileToMove)

	log.Println(expectedPath)

	if _, err := os.Stat(expectedPath); errors.Is(err, os.ErrNotExist) {
		test.Fail(err, expectedPath, t)
	}
}
