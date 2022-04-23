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
)

// Moves the file found at 'userspaceFile' to 'dotfilesDir' and creates a symlink in its original
// location pointing to it.
func AddFile(userspaceFile, dotfilesDir string) error {
	if err := checkIfFileExists(userspaceFile); err != nil {
		return err
	}

	if err := checkIfFileExists(dotfilesDir); err != nil {
		return err
	}

	// TODO:
	// make backup of userspace files
	// create path in dotfiles dir
	// copy files to dotfiles dir
	// remove files from userspace
	// create symlink in userspace

	return nil
}

// Backs up src to dst without modifying src. Returns path to dst.
// Both src and dst should be actual file paths, not directories.
func backupFileTemp(src, dst string) (string, error) {
	in, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("couldn't open src: %w", err)
	}
	defer in.Close()

	out, err := os.Create(dst)
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
		return "", fmt.Errorf("backup file was not found: %w", err)
	}

	return out.Name(), nil
}

// UpdateSymlinks walks over files and folders in the dotfiles dir, while updating their
// respective symlinks in the system relative to the placement in the dotfiles directory.
// `dotfilesDirPath` denotes the path to the dotfiles directory.
// `userSpacePath` denotes the root of where the symlinks can be found.
func UpdateSymlinks(dotfilesDirPath string, userSpacePath string) error {

	if err := checkIfFileExists(dotfilesDirPath); err != nil {
		return err
	}

	if err := checkIfFileExists(userSpacePath); err != nil {
		return err
	}

	absUserSpaceDir, err := filepath.Abs(userSpacePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %w", userSpacePath, err)
	}

	absDotfilesDirPath, err := filepath.Abs(dotfilesDirPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %w", dotfilesDirPath, err)
	}

	// Walkdir traverses the dotfiles dir with `p` denoting each file or directory in the dotfiles
	// directory and can be either a file or directory.
	return filepath.WalkDir(dotfilesDirPath, func(p string, d fs.DirEntry, err error) error {
		if p == dotfilesDirPath {
			return nil
		}

		// absolute path to each dotfile visited.
		absFilePath, err := filepath.Abs(p)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for %s: %w", p, err)
		}

		// construct the relative path to each file inside dotfiles dir by removing
		// the leading part of the path to the dotfiles dir.
		relativeDotfilesDirPath := strings.TrimPrefix(absFilePath, absDotfilesDirPath)

		// construct path to each expected loaction in user space imitating the
		// directory structure of the dotfiles directory.
		userFile := path.Join(absUserSpaceDir, relativeDotfilesDirPath)

		if IsFileSymlink(userFile) {
			err = UpdateSymlink(absFilePath, userFile)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateSymlink updates an existing symlink 'file' to point to 'pointTo'
func UpdateSymlink(pointTo string, file string) error {
	// symlink info: https://stackoverflow.com/questions/37345844/how-to-overwrite-a-symlink-in-go

	err := os.Remove(file)
	if err != nil {
		return fmt.Errorf("failed to unlink file: %s, %s, %+v", file, pointTo, err)
	}

	err = os.Symlink(pointTo, file)
	if err != nil {
		return fmt.Errorf("failed to create new symlink: %s, %s, %+v", file, pointTo, err)
	}

	fmt.Printf(Color("Updated: ", Green)+"%s -> %s.\n", file, pointTo)
	return nil
}

// IsFileSymlink returns true if the given path is an existsing symlink.
func IsFileSymlink(file string) bool {
	fileInfo, err := os.Lstat(file)
	if err != nil {
		fmt.Println(Color("Warning: ", Red), err)
		return false
	}

	return fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink
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
