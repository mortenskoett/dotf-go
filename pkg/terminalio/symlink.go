package terminalio

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mortenskoett/dotf-go/pkg/logging"
)

// UpdateSymlinks walks over files and folders in the dotfiles dir, while updating their respective
// symlinks in userspace relative to the placement in the dotfiles directory. If a matching symlink
// is not found in userspace, the file is ignored.
// `dotfilesDirPath` denotes the path to the dotfiles directory.
// `userSpacePath` denotes the root of where the symlinks can be found.
func UpdateSymlinks(userSpaceDir, dotfilesDir string) error {
	absUserSpaceDir, err := getAbsolutePath(userSpaceDir)
	if err != nil {
		return err
	}

	absDotfilesDir, err := getAbsolutePath(dotfilesDir)
	if err != nil {
		return err
	}

	// Walkdir traverses the dotfiles dir with `p` denoting each file or directory in the dotfiles
	// directory and can be either a file or directory.
	return filepath.WalkDir(dotfilesDir, func(p string, d fs.DirEntry, err error) error {
		if p == dotfilesDir {
			return nil
		}

		absFilePath, err := getAbsolutePath(p)
		if err != nil {
			return err
		}

		fileInUserspace, err := replacePrefixPath(absFilePath, absDotfilesDir, absUserSpaceDir)
		if err != nil {
			return err
		}

		exists, err := CheckIfFileExists(fileInUserspace)
		if err != nil {
			return err
		}
		if !exists {
			logging.Warn("Ignoring file because it doesn't exist in userspace: ", fileInUserspace)
			return nil
		}

		ok, err := IsFileSymlink(fileInUserspace)
		if err != nil {
			return err
		}

		if ok {
			err = updateSymlink(fileInUserspace, absFilePath)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// IsFileSymlink returns true if the given path is an existing symlink.
func IsFileSymlink(file string) (bool, error) {
	fileInfo, err := os.Lstat(file)
	if err != nil {
		return false, fmt.Errorf("failed to determine file is a symlink: %w", err)
	}
	return fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink, nil
}

// updateSymlink updates an existing symlink found at location 'fromDest' to point to an existing
// file 'toFile'. It is assumed that 'toFile' is a different location from what the symlink
// 'fromDest' previously pointed to. Otherwise the operation does not make sense.
// Example: A->B can be changed to A->C.
func updateSymlink(fromDest, toFile string) error {
	// symlink info: https://stackoverflow.com/questions/37345844/how-to-overwrite-a-symlink-in-go

	if err := deleteFile(fromDest); err != nil {
		return err
	}

	if err := createSymlink(fromDest, toFile); err != nil {
		return err
	}
	return nil
}

// Create a symlink at location 'symlinkDest' pointing to an actual file that should exist at
// 'fileSrc'.
func createSymlink(symlinkDest, fileSrc string) error {
	err := os.Symlink(fileSrc, symlinkDest)
	if err != nil {
		return fmt.Errorf("failed to create symlink from %s -> %s: %w", symlinkDest, fileSrc, err)
	}
	logging.Ok("Symlink successfully created from", symlinkDest, "->", fileSrc) // from symlink -> file
	return nil
}
