package main

import (
	"os"

	"github.com/mortenskoett/dotf-go/pkg/argparse"
	"github.com/mortenskoett/dotf-go/pkg/cli"
	"github.com/mortenskoett/dotf-go/pkg/logger"
)

func main() {
	// Parse input to command
	execName, cmdName, cliargs, err := argparse.ParseCliArguments(os.Args)
	if err != nil {
		switch err.(type) {
		case *argparse.ParseHelpFlagError:
			logger.LogSuccess(err)
		case *argparse.ParseNoArgumentError:
			logger.LogWarn(err)
		case *argparse.ParseInvalidArgumentError:
			logger.LogWarn(err)
		case *argparse.ParseError:
			logger.LogError(err)
		default:
			logger.LogError("unknown parser error:", err)
		}
		os.Exit(1)
	}

	// Create command
	cmd, err := cli.CreateCommand(execName, cmdName)
	if err != nil {
		switch err.(type) {
		case *cli.CmdUnknownCommand:
			logger.LogError(err)
		default:
			logger.LogError("unknown command error:", err)
		}
		os.Exit(1)
	}

	// If no errors run command
	err = cmd.Run(cliargs)
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
