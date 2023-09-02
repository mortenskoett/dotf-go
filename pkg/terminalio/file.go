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
	insideDotfiles bool   // Whether the filepath was to a file inside dotfiles
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
	if exists, _ := CheckIfFileExists(path); !exists {
		return "", &ErrFileNotFound{path}
	}

	return path, nil
}

// Checks if file exists by trying to open it. The given path should be absolute or relative to
// dotf executable. An error is returned if the file does not exist.
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

// Finds common prefix of two already split strings. E.g. '/a/b/c' and '/a/b/c/d' gives 'a/b/c'.
func FindCommonPathPrefix(path1, path2 string) (string, error) {
	return findCommonPath(path1, path2, commonPrefixFinder)
}

// Finds common suffix of two already split strings. E.g. '/a/b/c' and 'd/b/c' gives 'b/c'.
func FindCommonPathSuffix(path1, path2 string) (string, error) {
	return findCommonPath(path1, path2, commonSuffixFinder)
}

// Defines a strategy to analyze two slices of strings and built a slice from the input.
type sliceMerger func(long, short []string) []string

// Finds the common path overlap of two given paths using the given sliceMerger function. Strings
// are split using the default system specific delimeter.
func findCommonPath(path1, path2 string, mergerFunc sliceMerger) (string, error) {
	delimeter := filepath.Separator

	if path1 == "" || path2 == "" {
		return "", fmt.Errorf("given path was empty string")
	}

	if !strings.Contains(path1, string(delimeter)) || !strings.Contains(path2, string(delimeter)) {
		return "", fmt.Errorf("given input is not using '%v' as delimeter", string(delimeter))
	}

	if path1 == path2 {
		return path1, nil
	}

	// Split on specific runes to tokenize.
	splitFunc := func(r rune) bool {
		return r == delimeter
	}
	p1 := strings.FieldsFunc(path1, splitFunc)
	p2 := strings.FieldsFunc(path2, splitFunc)

	// Find longest slice in order to handle varying slice lengths.
	var commonstr []string
	if len(p1) >= len(p2) {
		commonstr = mergerFunc(p1, p2)
	} else {
		commonstr = mergerFunc(p2, p1)
	}

	// Reutrn common path as string.
	return filepath.Join(commonstr...), nil
}

func commonPrefixFinder(longest, shortest []string) []string {
	var sharedpath []string
	for i := 0; i < len(shortest); i++ {
		if longest[i] != shortest[i] {
			break
		}
		sharedpath = append(sharedpath, shortest[i])
	}
	return sharedpath
}

func commonSuffixFinder(longest, shortest []string) []string {
	// Loops backwards through slices potentially varying in length.
	// Example with indexes: short: 2,1,0 and long: 2+diff, 1+diff, 0+diff:
	// 0|1|2|3|5|6
	// 0|1|2

	var sharedpath []string
	var diff = len(longest) - len(shortest)

	for si := len(shortest) - 1; si >= 0; si-- {
		li := si + diff
		if longest[li] != shortest[si] {
			break
		}
		// Construct path backwards (right through left) requires preprending even if suboptimal
		sharedpath = append([]string{longest[li]}, sharedpath...)
	}
	return sharedpath
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

// Writes bytes to disk, overwriting file if it already exists.
func writeFile(fpath string, contents []byte) error {
	absPath, err := getAbsolutePath(fpath)
	if err != nil {
		return err
	}

	err = os.WriteFile(absPath, contents, 0644)
	if err != nil {
		return err
	}

	if exists, _ := CheckIfFileExists(absPath); !exists {
		return &ErrFileNotFound{absPath}
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
