package cli

import (
	"fmt"

//	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type moveCommand struct {}

func NewMoveCommand() *moveCommand {
	return &moveCommand{}
}

func (ma *moveCommand) Run(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("not enough arguments given")
	}

	//dotfilesDir := args[0]
	//symlinkRootDir := args[1]

	//err := terminalio.UpdateSymlinks(dotfilesDir, symlinkRootDir)
	//if err != nil {
	//	return fmt.Errorf("one or more errors happened while updating symlinks: ", err)
	//}

	//fmt.Println("\nAll symlinks updated successfully.")
	return nil
}

func (ma *moveCommand) Usage() CommandUsage { 
	return CommandUsage{
		Name: "move",
		Args: map[string]string {
			"from" : "Specifies dotfiles directory.",
			"to" : "Specifies userspace root directory where symlinks will be updated.",
		},
		Usage: "iterates through all files in 'from' and updates matching symlinks in 'to'.",
	}
}

func (ma *moveCommand) Description() string {
	return fmt.Sprintln(
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
}
