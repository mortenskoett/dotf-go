package main

import (
	"os"

	"github.com/mortenskoett/dotf-go/pkg/argparse"
	"github.com/mortenskoett/dotf-go/pkg/cli"
	"github.com/mortenskoett/dotf-go/pkg/logger"
)

func main() {
	// Parse input to command
	cmd, cliargs, err := argparse.HandleArguments(os.Args)
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

	// Run command
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
