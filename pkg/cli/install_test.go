package cli_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/mortenskoett/dotf-go/pkg/cli"
	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
	"github.com/mortenskoett/dotf-go/pkg/test"
)

// For tests.
type mockInteractor struct {
	b bytes.Buffer
}

func (s mockInteractor) ConfirmByUser(question string) bool {
	return cli.ConfirmByUser(question, &s.b)
}

func TestInstallExternalSymlink(t *testing.T) {
	// Arrange
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	externalDotfilesDir := env.BackupDir // We use this as external dotfiles dir.
	externalSharedFileToInstall := externalDotfilesDir.AddTempFile()
	externalSymlinkToInstall := externalDotfilesDir.CreateTempSymlink(externalSharedFileToInstall.Path)

	currentDotfilesDir := env.DotfilesDir
	currentUserspaceDir := env.UserspaceDir

	cliInput := &parsing.CommandlineInput{
		CommandName:    "install",
		PositionalArgs: []string{externalSymlinkToInstall.Path},
		Flags: parsing.NewFlagHolder(map[string]string{
			cli.FlagExternal: externalDotfilesDir.Path,
		}),
	}

	dotfConf := &parsing.DotfConfiguration{
		ConfigMetadata: &parsing.ConfigMetadata{},
		UserspaceDir:   currentUserspaceDir.Path,
		DotfilesDir:    currentDotfilesDir.Path,
	}

	// Use buffer instead of stdin
	var stdin bytes.Buffer
	stdin.Write([]byte("Y\n"))

	// Act
	cmd := cli.NewInstallCommand()
	cmd.UserInteractor = mockInteractor{b: stdin} // Insert buffer
	err := cmd.Run(cliInput, dotfConf)
	if err != nil {
		t.Errorf("%+v", err)
	}

	// Asserts

	// assert copied over config exists in dotfiles dir
	expectedDotfilesFile := filepath.Join(currentDotfilesDir.Path, externalSymlinkToInstall.Name)
	if exists, _ := terminalio.CheckIfFileExists(expectedDotfilesFile); !exists {
		t.Errorf("A new config pointed to from userspace should now exists now on in dotfiles on path: %v, opt error: %+v", expectedDotfilesFile, err)
	}

	// check if new file in dotfiles is in fact symlink
	if ok, err := terminalio.IsFileSymlink(expectedDotfilesFile); !ok || err != nil {
		t.Errorf("File in dotfiles dir should be a symlink pointing to the actual file (probably in shared). This is not a symlink: %s: %v", expectedDotfilesFile, err)
	}

	// assert symlink exists in userspace
	expectedUserspacePath := filepath.Join(currentUserspaceDir.Path, externalSymlinkToInstall.Name)
	if exists, err := terminalio.CheckIfFileExists(expectedUserspacePath); !exists {
		t.Errorf("A symlink pointing from userspace to dotfiles should exists now on path: %v, opt error: %+v", expectedUserspacePath, err)
	}

	// check if new file in userspace is in fact symlink
	if ok, err := terminalio.IsFileSymlink(expectedUserspacePath); !ok || err != nil {
		t.Errorf("File in userspace dir should be a symlink. This is not a symlink: %s: %v", expectedUserspacePath, err)
	}
}

func TestInstallExternalFile(t *testing.T) {
	// Arrange
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	externalDotfilesDir := env.BackupDir // We use this as external dotfiles dir.
	externalFileToInstall := externalDotfilesDir.AddTempFile()

	currentDotfilesDir := env.DotfilesDir
	currentUserspaceDir := env.UserspaceDir

	cliInput := &parsing.CommandlineInput{
		CommandName:    "install",
		PositionalArgs: []string{externalFileToInstall.Path},
		Flags: parsing.NewFlagHolder(map[string]string{
			cli.FlagExternal: externalDotfilesDir.Path,
		}),
	}

	dotfConf := &parsing.DotfConfiguration{
		ConfigMetadata: &parsing.ConfigMetadata{},
		UserspaceDir:   currentUserspaceDir.Path,
		DotfilesDir:    currentDotfilesDir.Path,
	}

	// Use buffer instead of stdin
	var stdin bytes.Buffer
	stdin.Write([]byte("Y\n"))

	// Act
	cmd := cli.NewInstallCommand()
	cmd.UserInteractor = mockInteractor{b: stdin} // Insert buffer
	err := cmd.Run(cliInput, dotfConf)
	if err != nil {
		t.Errorf("%+v", err)
	}

	// Asserts

	// assert copied over config exists in dotfiles dir
	expectedDotfilesFile := filepath.Join(currentDotfilesDir.Path, externalFileToInstall.Name)
	if exists, _ := terminalio.CheckIfFileExists(expectedDotfilesFile); !exists {
		t.Errorf("A new config pointed to from userspace should now exists now on in dotfiles on path: %v, opt error: %+v", expectedDotfilesFile, err)
	}

	// assert symlink exists in userspace
	expectedUserspacePath := filepath.Join(currentUserspaceDir.Path, externalFileToInstall.Name)
	if exists, err := terminalio.CheckIfFileExists(expectedUserspacePath); !exists {
		t.Errorf("A symlink pointing from userspace to dotfiles should exists now on path: %v, opt error: %+v", expectedUserspacePath, err)
	}

	// check if new file in userspace is in fact symlink
	if ok, err := terminalio.IsFileSymlink(expectedUserspacePath); !ok || err != nil {
		t.Errorf("File in userspace dir should be a symlink at %s: %v", expectedUserspacePath, err)
	}

}
