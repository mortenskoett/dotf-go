package main

import (
	"os"

	"github.com/mortenskoett/dotf-go/pkg/argparse"
	"github.com/mortenskoett/dotf-go/pkg/cli"
	"github.com/mortenskoett/dotf-go/pkg/logging"
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
	// Parse input to command
	cliargs, config, err := argparse.Parse(os.Args)
	if err != nil {
		switch err.(type) {
		case *argparse.ParseHelpFlagError:
			argparse.PrintFullHelp(cli.GetAvailableCommands(programName), programName, logo, programVersion)
			logging.Ok(err)
		case *argparse.ParseNoArgumentError:
			argparse.PrintBasicHelp(cli.GetAvailableCommands(programName), programName, logo, programVersion)
			logging.Warn(err)
		case *argparse.ParseInvalidArgumentError:
			logging.Warn(err)
		case *argparse.ParseConfigurationError:
			logging.Error(err)
		case *argparse.ParseError:
			logging.Error(err)
		default:
			logging.Error("unknown parser error:", err)
		}
		os.Exit(1)
	}

	// Create command
	cmd, err := cli.CreateCommand(programName, cliargs.CmdName)
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
