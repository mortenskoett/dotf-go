package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

const (
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

	err := terminalio.UpdateSymlinks(*dotfilesDir, *symlinkRootDir)
	if err != nil {
		log.Fatal("one or more errors happened while updating symlinks: ", err)
	}

	fmt.Println("\nAll symlinks updated successfully.")
}

func printDefaults() {
	fmt.Println(terminalio.Color(logo, terminalio.Blue))
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
