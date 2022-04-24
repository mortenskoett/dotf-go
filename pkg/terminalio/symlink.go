package terminalio

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mortenskoett/dotf-go/pkg/logger"
)

// Copies the file found at 'userspaceFile' to 'dotfilesDir' and creates a symlink in its original
// location pointing to it.
func AddFileCreateSymlink(userspaceFile, dotfilesDir string) error {
	absUserSpaceFile, err := getAbsolutePath(userspaceFile)
	if err != nil {
		return err
	}

	absUserSpaceFile, err = getAbsolutePath(dotfilesDir)
	if err != nil {
		return err
	}

	// TODO:
	// **OK** 	// make backup of userspace files
	// create path in dotfiles dir
	// copy files to dotfiles dir
	// remove files from userspace
	// create symlink in userspace

	_, err = backupFile(absUserSpaceFile)
	if err != nil {
		return err
	}

	return nil
}

// UpdateSymlinks walks over files and folders in the dotfiles dir, while updating their
// respective symlinks in the system relative to the placement in the dotfiles directory.
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

		// absolute path to each dotfile visited.
		absFilePath, err := filepath.Abs(p)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for %s: %w", p, err)
		}

		// construct the relative path to each file inside dotfiles dir by removing
		// the leading part of the path to the dotfiles dir.
		relativeDotfilesDirPath := strings.TrimPrefix(absFilePath, absDotfilesDir)

		// construct path to each expected loaction in user space imitating the
		// directory structure of the dotfiles directory.
		userFile := path.Join(absUserSpaceDir, relativeDotfilesDirPath)

		if isFileSymlink(userFile) {
			err = UpdateSymlink(absFilePath, userFile)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// TODO
// func create

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

	logger.LogSuccess("Updated:"+"%s -> %s.\n", file, pointTo)
	return nil
}
