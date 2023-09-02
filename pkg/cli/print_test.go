package cli_test

import (
	"bytes"
	"testing"

	"github.com/mortenskoett/dotf-go/pkg/cli"
)

func TestConfirmByUserYesReturnsTrue(t *testing.T) {
	testcases := []struct {
		input          string
		expectedResult bool
	}{
		{
			input:          "Y",
			expectedResult: true,
		},
		{
			input:          "yes",
			expectedResult: true,
		},
		{
			input:          "n",
			expectedResult: false,
		},
		{
			input:          "no",
			expectedResult: false,
		},
	}

	for _, tc := range testcases {
		var stdin bytes.Buffer
		stdin.Write([]byte(tc.input + "\n"))
		actualResult := cli.ConfirmByUser("Some great question?", &stdin)
		if tc.expectedResult != actualResult {
			t.Errorf("Handling of user input was not as expected: got: %+v, want %+v. Info: %+v", actualResult, tc.expectedResult, tc)
		}
	}
}
