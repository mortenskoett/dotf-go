package main

import (
	"os"

	"github.com/mortenskoett/dotf-go/pkg/argparse"
	"github.com/mortenskoett/dotf-go/pkg/logger"
)

func main() {
	cmd, cliargs, err := argparse.HandleArguments(os.Args)
	if err != nil {
		logger.LogFatal("fatal:", err)
	}

	err = cmd.Run()
	if err != nil {
		logger.LogFatal("fatal:", err)
	}

	logger.LogSuccess("OK:", cmd, cliargs)
}
