package cli

import (
	"fmt"

	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type moveCommand struct {}

func NewMoveCommand() *moveCommand {
	return &moveCommand{}
}

func (cmd *moveCommand) Run(args []string) error {
	if len(args) == 0 {
		printUsage(cmd)
		return nil
	}

	if len(args) == 1 && args[0] == "help" {
		printHelp(cmd)
		return nil
	}

	if len(args) != 2 {
		return fmt.Errorf("wrong number of arguments given")
	}

	dotfilesDir := args[0]
	symlinkRootDir := args[1]

	err := terminalio.UpdateSymlinks(dotfilesDir, symlinkRootDir)
	if err != nil {
		return err
	}

	fmt.Println("\nAll symlinks seems to have been updated successfully.")
	return nil
}

func (ma *moveCommand) Data() CommandData { 
	return CommandData{
		Name: "move",
		Args: map[string]string {
			"from" : "Specifies dotfiles directory.",
			"to" : "Specifies userspace root directory where symlinks will be updated.",
		},
		Desc: "Iterates through all files in 'from' and updates matching symlinks in 'to'.",
	}
}

func printUsage(c *moveCommand) {
	fmt.Println(BuildUsageText(c.Data()))
}

func printHelp(c *moveCommand) {
	printUsage(c)
	fmt.Println("")

	fmt.Println("")
	fmt.Print(terminalio.Color("Description:", terminalio.Yellow))
	fmt.Println(`
In case the dotfiles directory has been moved, it is necessary to update all symlinks
pointing back to the old location, to point to the new location. 
This application will iterate through all directories and files in 'from' and attempt to locate
a matching symlink in the same location relative to the given argument 'to'. The given path 'to' 
is the root of the user space, e.g. root of '~/' aka the home folder.`)

	fmt.Print(terminalio.Color("Notes:", terminalio.Yellow))
	fmt.Println(`
- It is expected that the Dotfiles directory has already been moved and that 'from' is the new location directory.
- User space describes where the symlinks are placed.
- Dotfiles directory is where the actual files are placed.
- The user space and the dotfiles directory must match in terms of file hierachy.
- Obs: currently if a symlink is not found in the user space, then it will not be touched, however a warning will be shown.`)
}
