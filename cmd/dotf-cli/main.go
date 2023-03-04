package main

import (
	"os"

	"github.com/mortenskoett/dotf-go/pkg/cli"
	"github.com/mortenskoett/dotf-go/pkg/logging"
	"github.com/mortenskoett/dotf-go/pkg/parsing"
)

const logo = `    _       _     __             _  _
 __| | ___ | |_  / _|  ___   __ | |(_)
/ _' |/ _ \|  _||  _| |___| / _|| || |
\__/_|\___/ \__||_|         \__||_||_|
`

const (
	programName string = "dotf-cli"
)

var programVersion = "" // Inserted by build process using -ldflags

func main() {
	run(os.Args)
}

func run(osargs []string) {
	// Parse cli args
	cliargs, err := parsing.ParseCommandlineInput(os.Args)
	if err != nil {
		handleParsingError(err)
	}

	// Parse dotf config
	config, err := parsing.ParseDotfConfig(cliargs.Flags)
	if err != nil {
		handleParsingError(err)
	}

	// Create command
	cmd, err := cli.CreateCommand(programName, cliargs.CommandName)
	if err != nil {
		switch err.(type) {
		case *cli.CmdUnknownCommand:
			logging.Error(err)
		default:
			logging.Error("unknown command init error:", err)
		}
		os.Exit(1)
	}

	// If no errors then run command
	err = cmd.Run(cliargs, config)
	if err != nil {
		switch err.(type) {
		case *cli.CmdHelpFlagError:
			logging.Ok(err)
		case *cli.CmdArgumentError:
			logging.Warn(err)
		case *cli.GitError:
			logging.Error(err)
		default:
			logging.Error("unknown run error:", err)
		}
		os.Exit(1)
	}

	logging.Ok("All good. Bye!")
}

// Handles and terminates program according to error type
func handleParsingError(err error) {
	if err != nil {
		switch err.(type) {
		case *parsing.ParseHelpFlagError:
			cli.PrintFullHelp(cli.GetAvailableCommands(programName), programName, logo, programVersion)
			logging.Ok(err)
		case *parsing.ParseNoArgumentError:
			cli.PrintBasicHelp(cli.GetAvailableCommands(programName), programName, logo, programVersion)
			logging.Ok(err)
		case *parsing.ParseInvalidArgumentError:
			logging.Warn(err)
		case *parsing.ParseConfigurationError:
			logging.Error(err)
		case *parsing.ParseError:
			logging.Error(err)
		default:
			logging.Error("unknown parser error:", err)
		}
		os.Exit(1)
	}
}
