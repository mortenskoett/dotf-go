package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	colorReset = "\033[0m"
	colorRed = "\033[31m"
	colorGreen = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue = "\033[34m"
	logo = `
	 ▄▄▄▄▄▄  ▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄    ▄▄   ▄▄ ▄▄▄▄▄▄▄ ▄▄   ▄▄ ▄▄▄▄▄▄▄ 
	█      ██       █       █       █  █  █▄█  █       █  █ █  █       █
	█  ▄    █   ▄   █▄     ▄█    ▄▄▄█  █       █   ▄   █  █▄█  █    ▄▄▄█
	█ █ █   █  █ █  █ █   █ █   █▄▄▄   █       █  █ █  █       █   █▄▄▄ 
	█ █▄█   █  █▄█  █ █   █ █    ▄▄▄█  █       █  █▄█  █       █    ▄▄▄█
	█       █       █ █   █ █   █      █ ██▄██ █       ██     ██   █▄▄▄ 
	█▄▄▄▄▄▄██▄▄▄▄▄▄▄█ █▄▄▄█ █▄▄▄█      █▄█   █▄█▄▄▄▄▄▄▄█ █▄▄▄█ █▄▄▄▄▄▄▄█
	`
)


func main() {
	dotfilesDir := flag.String("from", "", "Required. Specify dotfiles directory.")
	symlinkRootDir := flag.String("to", "", "Required. Specify user root directory where symlinks should be updated.")

	flag.Parse()

	if *dotfilesDir == "" || *symlinkRootDir == "" {
		printDefaults()
		os.Exit(0)
	}

	err := updateSymlinks(*dotfilesDir, *symlinkRootDir)
	if err != nil {
		log.Fatal("one or more errors happened while updating symlinks: ", err)
	}

	fmt.Println("\nAll symlinks updated successfully.")
}

// updateSymlinks walks over files and folders in the dotfiles dir, while updating their
// respective symlinks in the system relative to the placement in the dotfiles directory.
// `dotfilesDirPath` denotes the path to the dotfiles directory.
// `userSpacePath` denotes the root of where the symlinks can be found.
func updateSymlinks(dotfilesDirPath string, userSpacePath string) error {
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

		if isFileSymlink(userFile) {
			err = updateSymlink(absFilePath, userFile)
			if err != nil { 
				return err
			}
		}
		return nil
	})
}

// updateSymlink updates an existing symlink 'file' to point to 'pointTo'
func updateSymlink(pointTo string, file string) error {
// symlink info: https://stackoverflow.com/questions/37345844/how-to-overwrite-a-symlink-in-go

	err := os.Remove(file)
	if err != nil {
		return fmt.Errorf("failed to unlink file: %s, %s, %+v", file, pointTo, err)
	}

	err = os.Symlink(pointTo, file)
	if err != nil {
			return fmt.Errorf("failed to create new symlink: %s, %s, %+v", file, pointTo, err)
	}

	fmt.Printf(color("Updated: ", colorGreen) + "%s -> %s.\n", file, pointTo)
	return nil
}

// isFileSymlink returns true if the given path is an existsing symlink. 
func isFileSymlink(file string) bool {
	fileInfo, err := os.Lstat(file)
	if err != nil {
		fmt.Println(color("Warning: ", colorRed), err)
		return false
	}

	return fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink
}

func printDefaults() {
	fmt.Println(color(logo, colorBlue))
	fmt.Println(
`In case the dotfiles directory has been moved, it is necessary to update all symlinks
pointing back to the old location, to point to the new location. 
This application will iterate through all directories and files in 'from' and attempt to locate
a matching symlink in the same location relative to the given argument 'to'. The given path 'to' 
is the root of the user space, e.g. root of '~/' aka the home folder.

Notes:
- It is expected that the Dotfiles directory has already been moved and that 'from' is the new location directory.
- User space describes where the symlinks are placed.
- Dotfiles directory is where the actual files are placed.
- The user space and the dotfiles directory must match in terms of file hierachy.
- Obs: currently if a symlink is not found in the user space, then it will not be touched, however a warning
will be shown.`)

	fmt.Println("")
	fmt.Println("Usage:")
	flag.PrintDefaults()
}

func color(text string, color string) string {
	return string(color) + text + string(colorReset)
}
