package main

import (
	"os"

	"github.com/mortenskoett/dotf-go/pkg/argparse"
	"github.com/mortenskoett/dotf-go/pkg/cli"
	"github.com/mortenskoett/dotf-go/pkg/logger"
)

const logo = `    _       _     __             _  _
 __| | ___ | |_  / _|  ___   __ | |(_)
/ _' |/ _ \|  _||  _| |___| / _|| || |
\__/_|\___/ \__||_|         \__||_||_|
`

const (
	programName string = "dotf-cli"
)

func main() {
	// Parse input to command
	cliargs, config, err := argparse.Parse(os.Args)
	if err != nil {
		switch err.(type) {
		case *argparse.ParseHelpFlagError:
			argparse.PrintFullHelp(cli.GetAvailableCommands(programName), programName, logo)
			logger.LogSuccess(err)
		case *argparse.ParseNoArgumentError:
			argparse.PrintBasicHelp(cli.GetAvailableCommands(programName), programName, logo)
			logger.LogWarn(err)
		case *argparse.ParseInvalidArgumentError:
			logger.LogWarn(err)
		case *argparse.ParseConfigurationError:
			logger.LogError(err)
		case *argparse.ParseError:
			logger.LogError(err)
		default:
			logger.LogError("unknown parser error:", err)
		}
		os.Exit(1)
	}

	// Create command
	cmd, err := cli.CreateCommand(programName, cliargs.CmdName)
	if err != nil {
		switch err.(type) {
		case *cli.CmdUnknownCommand:
			logger.LogError(err)
		default:
			logger.LogError("unknown command error:", err)
		}
		os.Exit(1)
	}

	// If no errors then run command
	err = cmd.Run(cliargs, config)
	if err != nil {
		switch err.(type) {
		case *cli.CmdHelpFlagError:
			logger.LogSuccess(err)
		case *cli.CmdArgumentError:
			logger.LogWarn(err)
		default:
			logger.LogError("unknown run error:", err)
		}
		os.Exit(1)
	}
}
