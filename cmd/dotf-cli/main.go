// The main dotf application entry point.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

func init() {
	log.SetPrefix("dotf-cli: ")
}

const (
	logo = `
	 ▄▄▄▄▄▄  ▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄    ▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄ 
	█      ██       █       █       █  █       █       █
	█  ▄    █   ▄   █▄     ▄█    ▄▄▄█  █   ▄▄▄▄█   ▄   █
	█ █ █   █  █ █  █ █   █ █   █▄▄▄   █  █  ▄▄█  █ █  █
	█ █▄█   █  █▄█  █ █   █ █    ▄▄▄█  █  █ █  █  █▄█  █
	█       █       █ █   █ █   █      █  █▄▄█ █       █
	█▄▄▄▄▄▄██▄▄▄▄▄▄▄█ █▄▄▄█ █▄▄▄█      █▄▄▄▄▄▄▄█▄▄▄▄▄▄▄█
	`
)

func main() {
	args := os.Args[1:]

	if len(args) > 0 {
		parseArguments(args)
	} else {
		printHelp()
	}
}

func parseArguments(args []string) {
	switch command := args[0]; command {
	case "install":
		fmt.Println("install not implemented")
	case "add":
		fmt.Println("add not implemented")
	case "move":
		fmt.Println("move not implemented")
	}
}

func printHelp() {
	fmt.Println(terminalio.Color(logo, terminalio.Blue))
	fmt.Println(
`Usage:
dotf-cli <command> [possible args...]

Commands:
move ... ...
`)
}
