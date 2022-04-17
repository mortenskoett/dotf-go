package main

import (
	"os"

	"github.com/mortenskoett/dotf-go/pkg/argparse"
	"github.com/mortenskoett/dotf-go/pkg/command"
	"github.com/mortenskoett/dotf-go/pkg/logger"
)

func main() {
	cmd, cliargs, err := argparse.HandleArguments(os.Args)
	if err != nil {
		logger.LogFatal("fatal:", err)
	}

	err = cmd.Run(cliargs)

	switch err.(type) {
	case *command.CmdErrorSuccess:
		logger.LogSuccess("exiting:", err)
	default:
		logger.LogFatal("fatal: unknown error:", err)
	}
}
