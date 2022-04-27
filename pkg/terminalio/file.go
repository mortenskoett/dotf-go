package terminalio

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Copies the file found at 'userspaceFile' to 'dotfilesDir'.
func AddFileToDotfiles(userspaceFile, dotfilesDir string) error {
	absUserSpaceFile, err := getCheckAbsolutePath(userspaceFile)
	if err != nil {
		return err
	}

	absDotfilesDir, err := getCheckAbsolutePath(dotfilesDir)
	if err != nil {
		return err
	}

	log.Println("user", absUserSpaceFile, "dotf", absDotfilesDir)

	// TODO: Implement this function

	// _, err = backupFile(absUserSpaceFile)
	// if err != nil {
	// 	return err
	// }

	return nil
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

// Changes the leading path of 'filepath' from that of 'frompath' to that of 'topath'. It is assumed
// that 'filepath' points to a file that is contained in 'frompath'.
func changeLeadingPath(filepath, fromdir, todir string) (string, error) {
	relative, err := detachRelativePath(filepath, fromdir)
	if err != nil {
		return "", err
	}
	absTo, err := getCheckAbsolutePath(todir)
	if err != nil {
		return "", err
	}

	// Suffixes the relative path to that of the new location.
	newpath := path.Join(absTo, relative)
	return newpath, nil
}

// Detaches 'filepath' from 'basepath' and returns the path-suffix of 'filepath' which is relative
// to 'basepath'. It is assumed that basepath is part of filepath.
func detachRelativePath(filepath, basepath string) (string, error) {
	absFile, err := getCheckAbsolutePath(filepath)
	if err != nil {
		return "", err
	}
	absBase, err := getCheckAbsolutePath(basepath)
	if err != nil {
		return "", err
	}

	// Removes the leading part of 'absFile'. The part that matches that of absBase.
	relative := strings.TrimPrefix(absFile, absBase)
	return relative, nil
}

// Returns the absolute path. If the path does not point to anything an error is returned.
func getCheckAbsolutePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("cannot get absolute path of empty string")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to create absolute path for %s: %s", path, err)
	}

	if err := checkIfFileExists(absPath); err != nil {
		return "", fmt.Errorf("a file could not be found: %w", err)
	}

	return absPath, nil
}

// Checks if file exists by trying to open it. The given path should be absolute.
func checkIfFileExists(absPath string) error {
	_, err := os.Open(absPath)
	if errors.Is(err, os.ErrNotExist) {
		return &NotFoundError{absPath}
	}
	return nil
}
