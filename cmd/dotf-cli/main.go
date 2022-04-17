package main

import (
	"log"
	"os"

	"github.com/mortenskoett/dotf-go/pkg/argparse"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix(terminalio.Color("dotf-cli error: ", terminalio.Red))
	// argparse.HandleArguments(os.Args)

	log.Println(argparse.Parse(os.Args))
}
