package terminalio

import (
	"os"
	"testing"

	"github.com/mortenskoett/dotf-go/pkg/test"
)

func TestBackupFile(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	fileToBackup := env.UserspaceDir.AddTempFile().Name()
	expectedBackupPath := "/tmp/dotf-go/backups" + fileToBackup

	actual, err := backupFile(fileToBackup)
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(expectedBackupPath)

	if expectedBackupPath != actual {
		test.Fail(actual, expectedBackupPath, t)
	}
}
