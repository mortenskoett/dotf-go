package terminalio

import (
	"errors"
	"log"
	"os"
	"testing"
)

type testFilepathHandle struct {
	path string
}

type testEnvironment struct {
	dotfilesDir, userspaceDir, backupDir testFilepathHandle
	// Cleans up the environment. Should be called when done e.g. using defer
	cleanup func() error
}

func createTestDirHierachySetup() testEnvironment {
	dotfilesDir, err := os.MkdirTemp("", "dotfiles")
	if err != nil {
		log.Fatal(err)
	}

	userspaceDir, err := os.MkdirTemp("", "userspace")
	if err != nil {
		log.Fatal(err)
	}

	backupDir, err := os.MkdirTemp("", "backup")
	if err != nil {
		log.Fatal(err)
	}

	// Returns a func that should be called by the user of the environment after use
	cleanfunc := func(dotfilesdir, userspacedir, backupdir string) func() error {
		return func() error {
			if err := os.RemoveAll(dotfilesdir); err != nil {
				return err
			}
			if err := os.RemoveAll(userspacedir); err != nil {
				return err
			}
			if err := os.RemoveAll(backupdir); err != nil {
				return err
			}
			return nil
		}
	}

	return testEnvironment{
		dotfilesDir:  testFilepathHandle{dotfilesDir},
		userspaceDir: testFilepathHandle{userspaceDir},
		backupDir:    testFilepathHandle{backupDir},
		cleanup:      cleanfunc(dotfilesDir, userspaceDir, backupDir),
	}
}

func fail(actual, expected interface{}, t *testing.T) {
	t.Errorf("\nactual = %v\nexpected = %v", actual, expected)
}

func TestAddFileNotFoundError(t *testing.T) {
	userspacefile := "adsf"
	dotfilesdir := "adsf"

	expected := &NotFoundError{}
	actual := AddFile(userspacefile, dotfilesdir)

	if !errors.As(actual, &expected) {
		fail(actual, expected, t)
	}
}

func TestBackupFileTemp(t *testing.T) {
	env := createTestDirHierachySetup()
	defer env.cleanup()

	from := env.userspaceDir.path
	to := env.backupDir.path

	err := backupFileTemp(from, to)

	if err != nil {
		fail(err, "should not fail", t)
	}
}
