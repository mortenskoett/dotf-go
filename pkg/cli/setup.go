package cli

import (
	"fmt"

	"github.com/mortenskoett/dotf-go/pkg/logging"
	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type setupCommand struct {
	*commandBase
	UserInteractor UserInteractor
}

func NewSetupCommand() *setupCommand {
	name := "setup"
	overview := "Create a sensible default configuration."
	usage := name + " [--help]"
	args := []arg{}
	flags := []*parsing.Flag{}
	description := `
	Creates a configuration file needed for dotf-cli and dotf-tray to function. The file is created
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
		UserInteractor: StdInUserInteractor{},
	}
}

func (c *setupCommand) Run(args *parsing.CommandlineInput, _ *parsing.DotfConfiguration) error {
	config := parsing.NewSensibleConfiguration()

	cmap, err := parsing.ConvertConfigToMap(config)
	if err != nil {
		return fmt.Errorf("failed to create default config: %v", err)
	}

	bs := parsing.CreateSerializableConfig(cmap)

	err = terminalio.WriteFile(config.Filepath, bs, false)
	if err != nil {
		switch e := err.(type) {
		case *terminalio.ErrAbortOnOverwrite:
			logging.Warn(fmt.Sprintf("A config file already exists: %s", logging.Color(e.Path, logging.Green)))
			logging.Warn(logging.Color("Current configuration will be OVERWRITTEN if you say so", logging.Red))
			ok := c.UserInteractor.ConfirmByUser("Do you want to continue?")
			if ok {
				if err := terminalio.WriteFile(config.Filepath, bs, ok); err != nil {
					return err
				}
			} else {
				logging.Info("Aborted by user")
				return nil
			}
		default:
			return err
		}
	}
	logging.Ok("Configuration successfully created at", config.Filepath)

	return nil
}
