package cli

import (
	"fmt"

	"github.com/mortenskoett/dotf-go/pkg/parsing"
)

type setupCommand struct {
	*commandBase
}

func NewSetupCommand() *setupCommand {
	name := "setup"
	overview := "Create a sensible default configuration."
	usage := name + " [--help]"
	args := []arg{}
	flags := []*parsing.Flag{}
	description := `
	Creates a configuration file needed for dotf-cli and dotf-tray to function. The file is creatd
	at ~/.config/dotf/config and based on sensible defaults but please do go and check whether it is
	looking like expected. `

	return &setupCommand{
		commandBase: &commandBase{
			Name:        name,
			Overview:    overview,
			Usage:       usage,
			Args:        args,
			Flags:       flags,
			Description: description,
		},
	}
}

func (c *setupCommand) Run(args *parsing.CommandlineInput, conf *parsing.DotfConfiguration) error {
	fmt.Printf("configuration seems to exist already: %+v", *conf)
	fmt.Println("Got here done")
	return nil
}
