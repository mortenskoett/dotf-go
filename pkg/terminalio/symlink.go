package terminalio

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mortenskoett/dotf-go/pkg/logging"
)

// UpdateSymlinks walks over files and folders in the dotfiles dir, while updating their
// respective symlinks in userspace relative to the placement in the dotfiles directory.
// `dotfilesDirPath` denotes the path to the dotfiles directory.
// `userSpacePath` denotes the root of where the symlinks can be found.
func UpdateSymlinks(userSpaceDir, dotfilesDir string) error {
	absUserSpaceDir, err := GetAbsolutePath(userSpaceDir)
	if err != nil {
		return err
	}

	absDotfilesDir, err := GetAbsolutePath(dotfilesDir)
	if err != nil {
		return err
	}

	// Walkdir traverses the dotfiles dir with `p` denoting each file or directory in the dotfiles
	// directory and can be either a file or directory.
	return filepath.WalkDir(dotfilesDir, func(p string, d fs.DirEntry, err error) error {
		if p == dotfilesDir {
			return nil
		}

		absFilePath, err := GetAbsolutePath(p)
		if err != nil {
			return err
		}

		fileInUserspace, err := ChangeLeadingPath(absFilePath, absDotfilesDir, absUserSpaceDir)
		if err != nil {
			return err
		}

		if isFileSymlink(fileInUserspace) {
			err = UpdateSymlink(absFilePath, fileInUserspace)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateSymlink updates an existing symlink found at location 'fromDest' to point to an existing
// file 'toFile'
func UpdateSymlink(toFile string, fromDest string) error {
	// symlink info: https://stackoverflow.com/questions/37345844/how-to-overwrite-a-symlink-in-go

	if err := deleteFile(fromDest); err != nil {
		return err
	}

	if err := createSymlink(toFile, fromDest); err != nil {
		return err
	}
	return nil
}

// Create a symlink at location 'fromDest' pointing to an actual file that should exist at 'toFile'.
func createSymlink(toFile, fromDest string) error {
	err := os.Symlink(toFile, fromDest)
	if err != nil {
		return fmt.Errorf("failed to create symlink from %s -> %s: %w", toFile, fromDest, err)
	}
	logging.Ok("Symlink successfully created from", toFile, "->", fromDest) // from symlink -> file
	return nil
}

// isFileSymlink returns true if the given path is an existsing symlink.
func isFileSymlink(file string) bool {
	fileInfo, err := os.Lstat(file)
	if err != nil {
		logging.Warn("Warning:", err)
		return false
	}
	return fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink
}
