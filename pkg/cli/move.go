package cli

import (
	"fmt"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

const (
	commandName = "move"
)

// Implements Command interface
type moveCommand struct {
	programName string
}

func NewMoveCommand(programName string) *moveCommand {
	return &moveCommand{programName:programName}
}

func (c *moveCommand) Run(args []string) error {
	if len(args) == 0 {
		fmt.Println(GenerateUsage(c.programName, c))
		return nil
	}

	if args[len(args)-1] == "--help" {
		fmt.Println(GenerateUsage(c.programName, c))
		fmt.Print("Description:")
		fmt.Println(c.Description())
		return nil
	}

	if len(args) != 2 {
		return fmt.Errorf("wrong number of arguments given. Try adding --help.")
	}

	ok := confirmByUser("\nThis operation can be desctructive. Do you want to continue?")
	if !ok {
		fmt.Println("Aborted by user")
		return nil
	}

	dotfilesDir := args[0]
	symlinkRootDir := args[1]

	err := terminalio.UpdateSymlinks(dotfilesDir, symlinkRootDir)
	if err != nil {
		return err
	}

	fmt.Println("\nAll symlinks have been updated successfully.")
	return nil
}

func (c *moveCommand) Name() string {
	return commandName
}

func (c *moveCommand) Overview() string {
	return "Iterates through configs in 'dotfiles-dir' and updates matching symlinks in 'userspace-dir'."
}

func (c *moveCommand) Arguments() *[]Arg {
	return &[]Arg{
		{ Name: "dotfiles-dir", Description: "Path specifies re-located dotfiles directory."},
		{ Name: "userspace-dir", Description: "Specifies userspace root directory where symlinks will be updated."},
	}
}

func (c *moveCommand) Usage() string {
	return fmt.Sprintf("%s %s <dotfiles-dir> <userspace-dir> [--help]", c.programName, commandName)
}

func (c *moveCommand) Description() string {
	return `
	In case the dotfiles directory has been moved, it is necessary to update all symlinks pointing
	back to the old location, to point to the new location. This application will iterate through
	all directories and files in 'from' and attempt to locate a matching symlink in the same
	location relative to the given argument 'to'. The given path 'to' is the root of the user space,
	e.g. root of '~/' aka the home folder. Note that currently if a symlink is not found in the user
	space, then it will not be touched, however a warning will be shown.

	It is expected that the dotfiles directory has already been moved and that 'from' is the new
	location directory. 
	` }

