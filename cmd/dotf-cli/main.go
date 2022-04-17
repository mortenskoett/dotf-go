package main

import (
	"os"

	"github.com/mortenskoett/dotf-go/pkg/argparse"
	"github.com/mortenskoett/dotf-go/pkg/logger"
)

func main() {
	// argparse.HandleArguments(os.Args)

	logger.Log(argparse.Parse(os.Args))
}
