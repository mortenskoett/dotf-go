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

// Contains info about a files whereabouts in relation to dotfiles dir and userspace.
type fileLocationInfo struct {
	insideDotfiles bool   // Whether the org filepath was to a file inside dotfiles
	fileOrgPath    string // Absolute path of given file
	userspaceFile  string // Absolute path to file in userspace
	dotfilesFile   string // Absolute path to file in dotfiles
}

// Returns the absolute path from current directory or an error if the created path does not point
// to a file.
func GetAndValidateAbsolutePath(path string) (string, error) {
	path, err := getAbsolutePath(path)
	if err != nil {
		return "", err
	}
	if exists, _ := checkIfFileExists(path); !exists {
		return "", &FileNotFoundError{path}
	}

	return path, nil
}

// Returns a struct containing information about the given 'file' and its relations to dotfiles and
// to userspace. This is useful because often there are commands that should produce equal results
// both when called from dotfiles and userspace.
func getFileLocationInfo(file, userspaceDir, dotfilesDir string) (info *fileLocationInfo, err error) {
	info = &fileLocationInfo{}

	absFile, err := getAbsolutePath(file)
	if err != nil {
		return nil, err
	}

	absUserspaceDir, err := GetAndValidateAbsolutePath(userspaceDir)
	if err != nil {
		return nil, err
	}

	absDotfilesDir, err := GetAndValidateAbsolutePath(dotfilesDir)
	if err != nil {
		return nil, err
	}

	// Determine whether given filepath is inside or outside dotfiles dir
	if strings.HasPrefix(absFile, absDotfilesDir) {
		// Inside dotfiles
		userspaceFilepath, err := replacePrefixPath(absFile, absDotfilesDir, absUserspaceDir)
		if err != nil {
			return nil, err
		}

		info.userspaceFile = userspaceFilepath
		info.dotfilesFile = absFile
		info.insideDotfiles = true

	} else {
		// In userspace
		dotfilesFilepath, err := replacePrefixPath(absFile, absUserspaceDir, absDotfilesDir)
		if err != nil {
			return nil, err
		}

		info.userspaceFile = absFile
		info.dotfilesFile = dotfilesFilepath
		info.insideDotfiles = false
	}

	info.fileOrgPath = absFile

	return
}

// Writes bytes to disk overwriting file if on already exists.
func writeFile(fpath string, contents []byte) error {
	absPath, err := getAbsolutePath(fpath)
	if err != nil {
		return err
	}

	err = os.WriteFile(absPath, contents, 0644)
	if err != nil {
		return err
	}

	if exists, _ := checkIfFileExists(absPath); !exists {
		return &FileNotFoundError{absPath}
	}
	return nil
}

// Backs up file and returns the path to the backed up version of the file. The given path should be
// made absolute by the caller. The backed up file should not be expected to persist between reboots.
func backupFile(file string) (string, error) {
	const backupDir string = "/tmp/dotf-go/backups/"

	logging.Info("Creating backup")

	backupDst := filepath.Join(backupDir, file)
	path, err := copyFileOrDir(file, backupDst)
	if err != nil {
		return "", err
	}
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
	srcAbs, err := getAbsolutePath(src)
	if err != nil {
		return "", err
	}
	dstAbs, err := getAbsolutePath(dst)
	if err != nil {
		return "", err
	}

	// Copy all files recursively
	err = filepath.WalkDir(srcAbs, func(p string, d fs.DirEntry, err error) error {
		newfilepath, err := replacePrefixPath(p, srcAbs, dstAbs)

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

	logging.Ok("Directory successfully copied from", src, "->", dstAbs)

	return dstAbs, nil
}

// Copies src to dst without modifying src. Both src and dst should be actual file paths, not
// directories. The function uses absolute paths for both src and dst. Does not handle directories
// and will fail. The path of the new file is returned.
func copyFile(src, dst string) (string, error) {
	srcAbs, err := getAbsolutePath(src)
	if err != nil {
		return "", err
	}

	dstAbs, err := getAbsolutePath(dst)
	if err != nil {
		return "", err
	}

	fstat, err := os.Stat(srcAbs)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("the src file was not found: %w", err)
		}
		return "", fmt.Errorf("failed to stat src file: %w", err)
	}

	fperm := fstat.Mode().Perm()

	// Open src file
	fsrc, err := os.Open(srcAbs)
	if err != nil {
		return "", fmt.Errorf("failed to open src: %w", err)
	}
	defer fsrc.Close()

	// Create path to destination file
	dstPath := path.Dir(dstAbs)
	err = os.MkdirAll(dstPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create directories: %w", err)
	}

	// Create dst file
	fdst, err := os.Create(dstAbs)
	if err != nil {
		return "", fmt.Errorf("failed to create dst: %w", err)
	}
	defer fdst.Close()

	// Check dst file exists
	_, err = os.Stat(dstAbs)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("the created file was not found: %w", err)
		}
		return "", fmt.Errorf("failed to stat dst file: %w", err)
	}

	// Copy actual contents to new dst
	_, err = io.Copy(fdst, fsrc)
	if err != nil {
		return "", fmt.Errorf("failed to copy src to dst: %w", err)
	}

	// Set permissions from src file
	err = os.Chmod(dstAbs, fperm)
	if err != nil {
		return "", fmt.Errorf("failed to set permissions on dst file: %w", err)
	}

	logging.Ok("File successfully copied from", src, "->", dstAbs)

	// Return path to file
	return dstAbs, nil
}

// Replaces the shared prefix path in 'filepath' from that of 'fromdir' to that of 'todir'. It is
// assumed that 'filepath' points to a file that is contained under 'fromdir'.
// E.g. func("/a/b/c/d", "/a/b/", "/e/f/") -> "/e/f/c/d"
func replacePrefixPath(filepath, fromdir, todir string) (string, error) {
	relative, err := trimBasePath(filepath, fromdir)
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

// Trims away 'basepath' from 'filepath' returning the remaining suffix of 'filepath'. It is assumed
// that basepath is part of filepath. It is assumed that that 'basepath' is part of 'filepath'.
// Example: detach(dotfiles/, dotfiles/d1/d2/file.txt) -> d1/d2/file.txt
func trimBasePath(filepath, basepath string) (string, error) {
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

// Returns the absolute path from current directory.
func getAbsolutePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("cannot get absolute path of empty string")
	}

	path = expandTilde(path)

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to create absolute path for %s: %s", path, err)
	}

	return absPath, nil
}

// Expand ~/ to path of home dir of current user.
func expandTilde(path string) string {
	if strings.HasPrefix(path, "~/") {
		dirname, _ := os.UserHomeDir()
		return filepath.Join(dirname, path[2:])
	}
	return path
}

// Checks if file exists by trying to open it. The given path should be absolute or relative to
// dotf executable. An error is returned if the file does not exist.
func checkIfFileExists(absPath string) (bool, error) {
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
