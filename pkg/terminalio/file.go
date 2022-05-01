package terminalio

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mortenskoett/dotf-go/pkg/logger"
)

// Copies file in userspace to dotfiles dir using same relative path between 'homeDir' and
// 'dotfilesDir'. The file is backed up first.
func AddFileToDotfiles(userspaceFile, homeDir, dotfilesDir string) error {
	absUserspaceFile, err := getAndValidateAbsolutePath(userspaceFile)
	if err != nil {
		return err
	}

	absHomedir, err := getAndValidateAbsolutePath(homeDir)
	if err != nil {
		return err
	}

	absDotfilesDir, err := getAndValidateAbsolutePath(dotfilesDir)
	if err != nil {
		return err
	}

	// TODO: Implement this function

	// construct path relative to inside dotfiles dir
	// check if file already exists and exit early
	// create path in dotfiles dir
	// make backup of userspace files
	// copy files to dotfiles dir
	// remove files from userspace
	// create symlink in userspace

	// Create path inside dotfiles dir
	absNewFile, err := changeLeadingPath(absUserspaceFile, absHomedir, absDotfilesDir)
	if err != nil {
		return err
	}

	// Assert a file is not already in dotfiles dir at location
	exists, err := checkIfFileExists(absNewFile)
	if err != nil {
		return err
	}
	if exists {
		return &FileAlreadyExistsError{absNewFile}
	}

	// Backup file before copying it
	_, err = backupFile(userspaceFile)
	if err != nil {
		return err
	}

	// Copy file to dotfiles
	_, err = copyFile(userspaceFile, absNewFile)
	if err != nil {
		return err
	}

	// Remove file in userspace
	if err := deleteFile(absUserspaceFile); err != nil {
		return err
	}

	// Create symlink from userspace to the newly created file in dotfiles
	if err := createSymlink(absNewFile, absUserspaceFile); err != nil {
		return err
	}

	return nil
}

func deleteFile(file string) error {
	err := os.Remove(file)
	if err != nil {
		return fmt.Errorf("failed to delete file: %s: %w", file, err)
	}
	logger.Log("File successfully deleted at:", file)
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
	logger.Log("File successfully backed up from", file, "->", path)
	return path, nil
}

// Copies src to dst without modifying src. Both src and dst should be actual file paths, not
// directories. Returns path to dst. The function uses absolute paths for both src and dst.
func copyFile(src, dst string) (string, error) {
	srcAbs, err := getAbsolutePath(src)
	dstAbs, err := getAbsolutePath(dst)

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

// Changes the leading path of 'filepath' from that of 'fromdir' to that of 'todir'. It is assumed
// that 'filepath' points to a file that is contained in 'fromdir'.
func changeLeadingPath(filepath, fromdir, todir string) (string, error) {
	relative, err := detachRelativePath(filepath, fromdir)
	if err != nil {
		return "", err
	}
	absTo, err := getAbsolutePath(todir)
	if err != nil {
		return "", err
	}

	// Suffixes the relative path to that of the new location.
	newpath := path.Join(absTo, relative)
	return newpath, nil
}

// Detaches 'filepath' from 'basepath' and returns the path-suffix of 'filepath' which is relative
// to 'basepath'. It is assumed that basepath is part of filepath.
// Aka removes the prefix of filepath that matches basepath.
// Example:
// detach(dotfiles/, dotfiles/d1/d2/file.txt) -> d1/d2/file.txt
func detachRelativePath(filepath, basepath string) (string, error) {
	absFile, err := getAbsolutePath(filepath)
	if err != nil {
		return "", err
	}
	absBase, err := getAbsolutePath(basepath)
	if err != nil {
		return "", err
	}

	// Removes the leading part of 'absFile'. The part that matches that of absBase.
	relative := strings.TrimPrefix(absFile, absBase)
	return relative, nil
}

// Returns the absolute path from current directory or an error if the created path does not point
// to a file.
func getAndValidateAbsolutePath(path string) (string, error) {
	path, err := getAbsolutePath(path)
	if err != nil {
		return "", err
	}
	if exists, _ := checkIfFileExists(path); !exists {
		return "", &NotFoundError{path}
	}

	return path, nil
}

// Returns the absolute path from current directory.
func getAbsolutePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("cannot get absolute path of empty string")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to create absolute path for %s: %s", path, err)
	}

	return absPath, nil
}

// Checks if file exists by trying to open it. The given path should be absolute or relative to
// dotf executable. An error is return ed if the file does not exist.
func checkIfFileExists(absPath string) (bool, error) {
	_, err := os.Open(absPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// &NotFoundError{absPath}
			return false, nil
		}
		return false, err
	}
	return true, nil
}
