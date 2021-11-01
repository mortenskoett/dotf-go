// The main dotf application entry point.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mortenskoett/dotf-go/pkg/terminalio"
	"github.com/mortenskoett/dotf-go/pkg/cli"
)

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

var (
// commands contains the CLI commands that are currently implemented in dotf.
// The command is also the key.
	commands = map[string]cli.Command {
		"move": cli.NewMoveCommand(),
	}
)

func main() {
	log.SetFlags(0)
	log.SetPrefix(terminalio.Color("dotf-cli error: ", terminalio.Red))
	args := os.Args[1:]

	if len(args) > 0 {
		parseArguments(args)
	} else {
		printHelp()
	}
}

func getAction(input string) (cli.Command, error) {
	act, ok := commands[input]
	if ok {
		return act, nil
	}

	return nil, fmt.Errorf("%s command does not exist.", input)
}

func parseArguments(args []string) {
	cmd := args[0]
	action, err := getAction(cmd)
		if err != nil {
			printHelp()
			log.Fatal(err)
	}

	err = action.Run(args[1:])
	if err != nil {
		log.Fatal(err)
	}
}

func printHelp() {
	fmt.Println(terminalio.Color(logo, terminalio.Blue))
	fmt.Println("Usage:")
	fmt.Println("dotf-cli <command> [args...]\n")

	fmt.Println("command   <arg1>   <arg2>   description")

	// Print implemented commands.
	for _, c := range commands {
		cmd := c.Usage()
		fmt.Println(cmd.Name)
	}
}
