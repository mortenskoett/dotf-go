package cli

import (
	"fmt"
)

const (
	commandName = "move"
)

type moveCommand struct {
	programName string
}

func NewMoveCommand(programName string) *moveCommand {
	return &moveCommand{programName:programName}
}

func (c *moveCommand) Run(args []string) error {
	if len(args) == 0 {
		fmt.Println(c.Description())
		return nil
	}

	if len(args) == 1 && args[0] == "help" {
		fmt.Println(c.Usage())
		fmt.Println(c.Description())
		return nil
	}

	if len(args) != 2 {
		return fmt.Errorf("wrong number of arguments given")
	}

// TODO insert again
//	dotfilesDir := args[0]
//	symlinkRootDir := args[1]

//	err := terminalio.UpdateSymlinks(dotfilesDir, symlinkRootDir)
//	if err != nil {
//		return err
//	}

	fmt.Println("\nAll symlinks have been updated successfully.")
	return nil
}

func (c *moveCommand) Name() string {
	return c.programName + commandName
}

func (c *moveCommand) Overview() string {
	return "Iterates through all files in 'dotfiles-dir' and updates symlinks using identical dir structure starting at 'userspace-dir'."
}

func (c *moveCommand) Arguments() *[]Arg {
	return &[]Arg{
		Arg{ Name: "dotfiles-dir", Description: "Path specifies re-located dotfiles directory."},
		Arg{ Name: "userspace-dir", Description: "Specifies userspace root directory where symlinks will be updated."},
	}
}

func (c *moveCommand) Usage() string {
	return "hello"
}

func (c *moveCommand) Description() string {
	return `
In case the dotfiles directory has been moved, it is necessary to update all symlinks
pointing back to the old location, to point to the new location. 
This application will iterate through all directories and files in 'from' and attempt to locate
a matching symlink in the same location relative to the given argument 'to'. The given path 'to' 
is the root of the user space, e.g. root of '~/' aka the home folder.

- It is expected that the Dotfiles directory has already been moved and that 'from' is the new location directory.
- User space describes where the symlinks are placed.
- Dotfiles directory is where the actual files are placed.
- The user space and the dotfiles directory must match in terms of file hierachy.
- Obs: currently if a symlink is not found in the user space, then it will not be touched, however a warning will be shown.`
}
