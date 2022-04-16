package cli

import "fmt"

type addCommand struct {
	CommandBase
}

func NewAddCommand(programName, commandName string) *addCommand {
	return &addCommand{
		CommandBase{
			programName: programName,
			commandName: commandName}}
}

func (c *addCommand) Run(args []string) error {
	// checkCmdArguments(args, c)
	if len(args) != 2 {
		return fmt.Errorf("wrong number of arguments given. Try adding --help.")
	}

	ok := confirmByUser("\nThis operation can be desctructive. Do you want to continue?")
	if !ok {
		fmt.Println("Aborted by user")
		return nil
	}

	// Actual operation
	return nil
}

func (c *addCommand) CmdName() string {
	return c.commandName
}

func (c *addCommand) Overview() string {
	return "to impl overview"
}

func (c *addCommand) Arguments() *[]Arg {
	return &[]Arg{
		{Name: "to impl name", Description: "to impl desc"},
		{Name: "to impl name", Description: "to impl desc"},
	}
}

func (c *addCommand) Usage() string {
	return fmt.Sprintf("%s %s <dotfiles-dir> <userspace-dir> [--help]", c.programName, c.commandName)
}

func (c *addCommand) Description() string {
	return ""
}

func (c *addCommand) ProgName() string {
	return c.programName
}
