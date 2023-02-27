package terminalio

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/mortenskoett/dotf-go/pkg/test"
)

func Test_UpdateSymlinks_recursively_changes_symlinks_in_userspace(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	// Construct dotfiles dir
	dotfiles := env.DotfilesDir
	dsomedir := dotfiles.AddTempDir("mydir1/mydir2")
	dsomefile := dsomedir.AddTempFile()

	// Construct userspace with symlinks with same paths/names as files in dotfiles
	userspace := env.UserspaceDir
	usomedir := userspace.AddTempDir("mydir1/mydir2")
	usomefile := userspace.AddTempFile()
	somesymlinkpath := filepath.Join(usomedir.Path, filepath.Base(dsomefile.Path))

	// Create symlink with same path as somefile in dotfiles to random file in userspace
	err := createSymlink(somesymlinkpath, usomefile.Path)
	if err != nil {
		t.Fail()
	}

	// Assert symlink exists
	ok, err := isFileSymlink(somesymlinkpath)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		test.FailHardMsg("This file should be a symlink", ok, true, t)
	}

	pathToLinkedFile, err := filepath.EvalSymlinks(somesymlinkpath)
	if err != nil {
		t.Fatal(err)
	}
	if pathToLinkedFile != usomefile.Path {
		test.FailHardMsg("This symlink should point to file in userspace", pathToLinkedFile, usomefile, t)
	}

	err = UpdateSymlinks(userspace.Path, dotfiles.Path)
	if err != nil {
		test.Fail(err, "Should not fail", t)
	}

	pathToLinkedFileAfterChange, err := filepath.EvalSymlinks(somesymlinkpath)
	if err != nil {
		t.Fatal(err)
	}
	if pathToLinkedFileAfterChange != dsomefile.Path {
		test.FailHardMsg("This symlink should point to file in userspace", pathToLinkedFile, usomefile, t)
	}
}

func Test_updateSymlink_makes_symlink_point_to_different_location(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	somefile := env.DotfilesDir.AddTempFile().Path
	anotherfile := env.DotfilesDir.AddTempFile().Path

	symlinkToSomefile := env.DotfilesDir.CreateTempSymlink(somefile).Path

	type testinput struct {
		from       string
		to         string
		shouldfail bool
	}

	testcases := []testinput{
		{
			from:       symlinkToSomefile,
			to:         anotherfile,
			shouldfail: false,
		},
		{
			from:       symlinkToSomefile,
			to:         "",
			shouldfail: true,
		},
		{
			from:       "",
			to:         anotherfile,
			shouldfail: true,
		},
	}

	for _, in := range testcases {
		err := updateSymlink(in.from, in.to)
		if err != nil {
			if !in.shouldfail {
				test.Fail(err, "updateSymlink should have updated path", t)
			}
		}
	}
}

func Test_createSymlink_fails_with_invalid_paths(t *testing.T) {
	type input struct {
		fpath      string
		symlink    string
		shouldfail bool
	}

	env := test.NewTestEnvironment()
	defer env.Cleanup()
	somefile := env.DotfilesDir.AddTempFile().Path

	sadpaths := []input{
		{
			fpath:      "",
			symlink:    fmt.Sprintf("%s-%s", somefile, "sym1"),
			shouldfail: true,
		},
		{
			fpath:      somefile,
			symlink:    "",
			shouldfail: true,
		},
		{
			fpath:      "",
			symlink:    "",
			shouldfail: true,
		},
	}

	for _, in := range sadpaths {
		err := createSymlink(in.symlink, in.fpath)

		// Should not be nil!
		if err == nil {
			test.Fail(in, "Test didn't fail as expected!", t)
		}
	}
}

func Test_createSymlink_creates_dir_file_symlink(t *testing.T) {
	type input struct {
		fpath      string
		symlink    string
		shouldfail bool
	}

	env := test.NewTestEnvironment()
	defer env.Cleanup()
	somefile := env.DotfilesDir.AddTempFile().Path
	somedir := env.DotfilesDir.AddTempDir("mylittledir").Path

	happypaths := []input{
		{
			fpath:   somefile,
			symlink: fmt.Sprintf("%s-%s", somefile, "sym1"),
		},
		{
			fpath:   somedir,
			symlink: fmt.Sprintf("%s-%s", somedir, "sym2"),
		},
	}

	for _, in := range happypaths {
		err := createSymlink(in.symlink, in.fpath)
		if err != nil {
			test.FailHard(err, "shouldn't fail here", t)
		}

		ok, err := isFileSymlink(in.symlink)
		if err != nil {
			test.FailHard(err, "Not expected to fail", t)
		}
		if !ok {
			test.FailMsg("isFileSymlink should return true here", ok, true, t)
		}
	}
}

func Test_isFileSymlink_determines_correctly_a_symlink(t *testing.T) {
	type input struct {
		symlink   string
		want      bool
		shoulderr bool
	}

	env := test.NewTestEnvironment()
	defer env.Cleanup()

	file := env.DotfilesDir.AddTempFile()
	symlink := env.UserspaceDir.CreateTempSymlink(file.Path)

	tests := []input{
		{
			// handle symlink
			symlink:   symlink.Path,
			want:      true,
			shoulderr: false,
		},
		{
			// handle non symlink
			symlink:   file.Path,
			want:      false,
			shoulderr: false,
		},

		{
			// handle error case
			symlink:   "yadda",
			want:      false,
			shoulderr: true, // here
		},
	}

	for _, tc := range tests {
		ok, err := isFileSymlink(tc.symlink)
		if err != nil {
			if !tc.shoulderr {
				test.Fail(err, "shouldn't fail here", t)
			}
		}
		if ok != tc.want {
			test.FailMsg("isFileSymlink error", ok, tc.want, t)
		}
	}
}
