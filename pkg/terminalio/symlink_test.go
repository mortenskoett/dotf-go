package terminalio

import (
	"fmt"
	"testing"

	"github.com/mortenskoett/dotf-go/pkg/test"
)

func Test_createSymlink_fails_with_invalid_paths(t *testing.T) {
	type input struct {
		fpath      string
		symlink    string
		shouldfail bool
	}

	env := test.NewTestEnvironment()
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
	somefile := env.DotfilesDir.AddTempFile().Path
	somedir := env.DotfilesDir.AddTempDir("mylittledir").Path

	happypaths := []input{
		{
			fpath:      somefile,
			symlink:    fmt.Sprintf("%s-%s", somefile, "sym1"),
		},
		{
			fpath:      somedir,
			symlink:    fmt.Sprintf("%s-%s", somedir, "sym2"),
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
