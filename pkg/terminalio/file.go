package terminalio

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mortenskoett/dotf-go/pkg/logging"
)

// Reverts the insertion of a file into the dotfiles directory and return it to its original
// location in userspace. The symlink is removed first. The operation can be applied both to the
// symlink in userspace and the actual file in the dotfiles directory.
func RevertDotfile(file, homeDir, dotfilesDir string) error {

	// TODO
	// a) If filepath is in userspace
	//	if NOT a symlink then return
	//	verify file exists at relative location in dotfiles

	// b) If filepath is in dotfiles
	//	if a symlink then return
	//	verify symlink exists in userspace using same relative path
	// if strings.HasPrefix(absFile, absDotfilesDir) {
	// 	logging.Ok("Is in dotfiles!")
	// }

	// both)
	//	remove symlink
	//	backup file in dotfiles
	//	copy file to userspace
	//	delete file in dotfiles

	absFile, err := GetAndValidateAbsolutePath(file)
	if err != nil {
		return err
	}

	absHomedir, err := GetAndValidateAbsolutePath(homeDir)
	if err != nil {
		return err
	}

	absDotfilesDir, err := GetAndValidateAbsolutePath(dotfilesDir)
	if err != nil {
		return err
	}

	dhome, _ := detachRelativePath(absFile, absHomedir)
	ddot, _ := detachRelativePath(absFile, absDotfilesDir)
	yeah, _ := detachRelativePath(absDotfilesDir, absHomedir)
	logging.Log("detachhome", dhome)
	logging.Log("detachdotf", ddot)
	logging.Log("yeah", yeah)
	logging.Log("file", absFile)
	logging.Log("homedir", absHomedir)
	logging.Log("dotfiles", absDotfilesDir)

	// If dotfiles dir is contained inside homedir
	// if strings.HasPrefix(absDotfilesDir, absHomedir) {
	// 	logging.Ok("Dotfiles dir is placed under homedir")

	// 	// Then we can know that
	// 	if strings.HasPrefix(absFile, absDotfilesDir) {
	// 		logging.Ok("Is in dotfiles!")
	// 	} else {
	// 		logging.Ok("Is in userspace!")
	// 	}
	// }

	// if strings.HasPrefix(absFile, absHomedir) {
	// 	logging.Ok("Is in userspace!")
	// }

	return nil
}

// Copies file in userspace to dotfiles dir using same relative path between 'homeDir' and
// 'dotfilesDir'. The file is backed up first.
func AddFileToDotfiles(userspaceFile, homeDir, dotfilesDir string) error {
	absUserspaceFile, err := GetAndValidateAbsolutePath(userspaceFile)
	if err != nil {

		return err
	}

	absHomedir, err := GetAndValidateAbsolutePath(homeDir)
	if err != nil {
		return err
	}

	absDotfilesDir, err := GetAndValidateAbsolutePath(dotfilesDir)
	if err != nil {
		return err
	}

	// Create path inside dotfiles dir
	absNewDotFile, err := ChangeLeadingPath(absUserspaceFile, absHomedir, absDotfilesDir)
	if err != nil {
		return err
	}

	// Assert a file is not already in dotfiles dir at location
	exists, err := CheckIfFileExists(absNewDotFile)
	if err != nil {
		return err
	}
	if exists {
		return &FileAlreadyExistsError{absNewDotFile}
	}

	// Backup file before copying it
	_, err = BackupFile(absUserspaceFile)
	if err != nil {
		return err
	}

	// Copy file to dotfiles
	_, err = copyFileOrDir(absUserspaceFile, absNewDotFile)
	if err != nil {
		return err
	}

	// Remove file in userspace
	if err := deleteFileOrDir(absUserspaceFile); err != nil {
		return err
	}

	// Create symlink from userspace to the newly created file in dotfiles
	if err := createSymlink(absNewDotFile, absUserspaceFile); err != nil {
		return err
	}

	return nil
}

// Backs up file and returns the path to the backed up version of the file. The given path should be
// made absolute by the caller. The backed up file should not be expected to persist between reboots.
func BackupFile(file string) (string, error) {
	const backupDir string = "/tmp/dotf-go/backups/"

	backupDst := backupDir + file
	path, err := copyFileOrDir(file, backupDst)
	if err != nil {
		return "", err
	}
	logging.Ok("Successfully created backup from", file, "->", path)
	return path, nil
}

// Will determine whether given 'src' points to a file or a directory and handle it accordingly. The
// function copies src to dst without modifying src. Src should be either a file or directory and
// dst should be a file path. Will copy directories recursively. Returns path to dst.
func copyFileOrDir(src, dst string) (string, error) {
	isDir, err := isDirectory(src)
	if err != nil {
		return "", fmt.Errorf("failed to determine if file was a directory: %w", err)
	}

	if isDir {
		return copyDir(src, dst)
	}

	return copyFile(src, dst)
}

// Copies a directory and its contents recursively from src to dst and return the absolute path to
// dst.
func copyDir(src, dst string) (string, error) {
	srcAbs, err := GetAbsolutePath(src)
	if err != nil {
		return "", err
	}
	dstAbs, err := GetAbsolutePath(dst)
	if err != nil {
		return "", err
	}

	// Copy all files recursively
	err = filepath.WalkDir(srcAbs, func(p string, d fs.DirEntry, err error) error {
		newfilepath, err := ChangeLeadingPath(p, srcAbs, dstAbs)

		isDir, err := isDirectory(p)
		if err != nil {
			return err
		}
		if isDir {
			return os.MkdirAll(newfilepath, os.ModePerm)
		}

		_, err = copyFile(p, newfilepath)
		if err != nil {
			return err
		}

		return nil
	})

	return dstAbs, nil
}

// Copies src to dst without modifying src. Both src and dst should be actual file paths, not
// directories. The function uses absolute paths for both src and dst. Does not handle directories
// and will fail.
func copyFile(src, dst string) (string, error) {
	srcAbs, err := GetAbsolutePath(src)
	if err != nil {
		return "", err
	}
	dstAbs, err := GetAbsolutePath(dst)
	if err != nil {
		return "", err
	}

	in, err := os.Open(srcAbs)
	if err != nil {
		return "", fmt.Errorf("couldn't open src: %w", err)
	}
	defer in.Close()

	// Create path to destination file
	dstPath := path.Dir(dstAbs)
	err = os.MkdirAll(dstPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("couldn't because of nested directores in dst: %w", err)
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
func ChangeLeadingPath(filepath, fromdir, todir string) (string, error) {
	relative, err := detachRelativePath(filepath, fromdir)
	if err != nil {
		return "", err
	}
	absTo, err := GetAbsolutePath(todir)
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
	absFile, err := GetAbsolutePath(filepath)
	if err != nil {
		return "", err
	}
	absBase, err := GetAbsolutePath(basepath)
	if err != nil {
		return "", err
	}

	// Removes the leading part of 'absFile'. The part that matches that of absBase.
	relative := strings.TrimPrefix(absFile, absBase)
	return relative, nil
}

// Returns the absolute path from current directory or an error if the created path does not point
// to a file.
func GetAndValidateAbsolutePath(path string) (string, error) {
	path, err := GetAbsolutePath(path)
	if err != nil {
		return "", err
	}
	if exists, _ := CheckIfFileExists(path); !exists {
		return "", &NotFoundError{path}
	}

	return path, nil
}

// Returns the absolute path from current directory.
func GetAbsolutePath(path string) (string, error) {
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
func CheckIfFileExists(absPath string) (bool, error) {
	_, err := os.Open(absPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func isDirectory(src string) (bool, error) {
	in, err := os.Open(src)
	if err != nil {
		return false, fmt.Errorf("couldn't open src: %w", err)
	}
	defer in.Close()

	file, err := in.Stat()
	if err != nil {
		return false, fmt.Errorf("couldn't stat src: %w", err)
	}

	if file.IsDir() {
		return true, nil
	}

	return false, nil
}

func deleteFileOrDir(path string) error {
	isDir, err := isDirectory(path)
	if err != nil {
		return err
	}

	if isDir {
		return deleteDirectory(path)
	}

	return deleteFile(path)
}

func deleteDirectory(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("failed to delete directory: %s: %w", path, err)
	}

	logging.Ok("Directory successfully deleted at", path)
	return nil
}

func deleteFile(file string) error {
	err := os.Remove(file)
	if err != nil {
		return fmt.Errorf("failed to delete file: %s: %w", file, err)
	}
	logging.Ok("File successfully deleted at", file)
	return nil
}
