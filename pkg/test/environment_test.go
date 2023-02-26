package test_test

import (
	"os"
	"testing"

	"github.com/mortenskoett/dotf-go/pkg/test"
)

func Test_AddTempDir_adds_dir_to_env(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	paths := []string{
		env.DotfilesDir.AddTempDir("testdotfiles").Path,
		env.UserspaceDir.AddTempDir("testuserspace").Path,
		env.BackupDir.AddTempDir("testbackup").Path,
	}

	for _, l := range paths {
		_, err := os.Lstat(l)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func Test_AddTempFile_adds_file_to_env(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	paths := []string{
		env.DotfilesDir.AddTempFile().Path,
		env.UserspaceDir.AddTempFile().Path,
		env.BackupDir.AddTempFile().Path,
	}

	for _, l := range paths {
		_, err := os.Lstat(l)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func Test_CreateTempSymlink_creates_symlink(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	dotfiles := env.DotfilesDir
	dfile := dotfiles.AddTempFile()

	userspace := env.UserspaceDir
	symlinkpath := userspace.CreateTempSymlink(dfile.Path)

	fileInfo, err := os.Lstat(symlinkpath)
	if err != nil {
		t.Fatal(err)
	}

	if !(fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink) {
		test.FailMsg("expected path to be symlink", fileInfo.Name, "a symlink", t)
	}
}
