package parsing_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/test"
)

func TestCommandlineInputParsing(t *testing.T) {
	type testinput struct {
		what       string
		args       []string
		shouldfail bool
		want       *parsing.CommandLineInput
	}

	testcases := []testinput{
		{
			what:       "command name is parsed ok",
			args:       []string{"executable", "command"},
			shouldfail: true, // ParseNoArgumentError
			want: &parsing.CommandLineInput{
				CommandName:    "command",
				PositionalArgs: []string{},
				Flags: &parsing.CommandLineFlags{
					ValueFlags: map[string]string{},
					BoolFlags:  map[string]bool{},
				},
			},
		},
		{
			what:       "full happy example ok",
			args:       []string{"executable", "command", "arg1", "arg2", "--valueflag1", "value1", "--boolflag1"},
			shouldfail: false,
			want: &parsing.CommandLineInput{
				CommandName:    "command",
				PositionalArgs: []string{"arg1", "arg2"},
				Flags: &parsing.CommandLineFlags{
					ValueFlags: map[string]string{"valueflag1": "value1"},
					BoolFlags:  map[string]bool{"boolflag1": true},
				},
			},
		},
		{
			what: "double value no flag name fails",
			args: []string{
				"executable",
				"command",
				"--valueflag1",
				"value1err",
				"value1",
			},
			shouldfail: true,
			want:       nil,
		},
	}

	for _, tc := range testcases {
		actual, err := parsing.ParseCommandlineInput(tc.args)
		if err != nil && !tc.shouldfail {
			t.Fatal("FAIL", err)
		}

		// Assert on error
		if err != nil && tc.shouldfail {
			switch err.(type) {
			case *parsing.ParseNoArgumentError:
				// ignore because it is used as signal for frontends
			case *parsing.ParseInvalidFlagError:
				// ignore because it is used as signal for frontends
			default:
				t.Fatal("FAIL", err)
			}
		}

		diff := cmp.Diff(tc.want, actual)
		if diff != "" {
			t.Errorf("failed to parse command line args for test: %s\n%s", tc.what, diff)
			test.PrintJSON("Actual", actual)
			test.PrintJSON("Want", tc.want)
		}
		t.Logf("PASS %s", tc.what)
	}
}
