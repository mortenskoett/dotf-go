package test_test

import (
	"os"
	"testing"

	"github.com/mortenskoett/dotf-go/pkg/test"
)

func Test_CreateSymlink_creates_symlink(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	dotfiles := env.DotfilesDir
	dfile := dotfiles.AddTempFile()

	userspace := env.UserspaceDir
	symlinkpath := userspace.CreateSymlink(dfile.Name())

	fileInfo, err := os.Lstat(symlinkpath)
	if err != nil {
		t.Fatal(err)
	}

	if !(fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink) {
		test.FailMsg("expected path to be symlink", fileInfo.Name, "a symlink", t)
	}
}
