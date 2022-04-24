package terminalio

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/mortenskoett/dotf-go/pkg/logger"
)

// isFileSymlink returns true if the given path is an existsing symlink.
func isFileSymlink(file string) bool {
	fileInfo, err := os.Lstat(file)
	if err != nil {
		logger.LogWarn("Warning:", err)
		return false
	}
	return fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink
}

// Backs up file and returns the path to the backed up version of the file. The backed up file
// should not be expected to persist between reboots.
func backupFile(file string) (string, error) {
	const backupDir string = "/tmp/dotf-go/backups/"

	backupDst := backupDir + file
	path, err := copyFile(file, backupDst)
	if err != nil {
		return "", err
	}
	return path, nil
}

// Copies src to dst without modifying src. Both src and dst should be actual file paths, not
// directories. Returns path to dst. The function uses absolute paths for both src and dst.
func copyFile(src, dst string) (string, error) {
	srcAbs, err := filepath.Abs(src)
	dstAbs, err := filepath.Abs(dst)

	in, err := os.Open(srcAbs)
	if err != nil {
		return "", fmt.Errorf("couldn't open src: %w", err)
	}
	defer in.Close()

	// Create path to destination file
	dstPath := path.Dir(dstAbs)
	err = os.MkdirAll(dstPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("couldn't nested directores in dst: %w", err)
	}

	out, err := os.Create(dstAbs)
	if err != nil {
		return "", fmt.Errorf("couldn't create dst: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return "", fmt.Errorf("couldn't copy src to dst: %w", err)
	}

	_, err = os.Stat(out.Name())
	if errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("the created file was not found: %w", err)
	}

	// Return path to file
	return out.Name(), nil
}

func getAbsolutePath(path string) (string, error) {
	if err := checkIfFileExists(path); err != nil {
		return "", fmt.Errorf("failed to get absolute path for %s: %w", path, err)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for %s: %s", path, err)
	}
	return absPath, nil
}

func checkIfFileExists(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		return &NotFoundError{absPath}
	}

	return nil
}
