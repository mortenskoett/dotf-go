package terminalio

import (
	"testing"

	"github.com/mortenskoett/dotf-go/pkg/test"
)

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
