// The main dotf-go application entry point.
package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

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
	commands = map[string]cli.Command {
		"move": cli.NewMoveCommand(),
	}
)

func main() {
	log.SetFlags(0)
	log.SetPrefix(terminalio.Color("dotf-cli error: ", terminalio.Red))
	args := os.Args[1:]

	if len(args) > 0 {
		handleArguments(args)
	} else {
		printHelp()
	}
}

func getCommand(input string) (cli.Command, error) {
	cmd, ok := commands[input]
	if ok {
		return cmd, nil
	}

	return nil, fmt.Errorf("%s command does not exist.", input)
}

func handleArguments(args []string) {
	input := args[0]

	if input == "help" {
		printHelp()
		return
	}

	cmd, err := getCommand(input)
		if err != nil {
			log.Fatal(err)
	}

	err = cmd.Run(args[1:])
	if err != nil {
		log.Fatal(err)
	}
}

func printHelp() {
	fmt.Println(terminalio.Color(logo, terminalio.Blue))
	fmt.Println(terminalio.Color("Usage: dotf-go <command> [args]", terminalio.Yellow))
	fmt.Println(terminalio.Color("Usage: dotf-go <command> help", terminalio.Yellow))
	fmt.Println("")

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 4, ' ', 0)

	fmt.Fprintln(w, "COMMAND\tARGS\tDESCRIPTION")

	// Print commands.
	for _, c := range commands {
		cmdata := c.Data()

		var arguments string
		for arg, _ := range cmdata.Args {
			arguments += "<" + arg + ">" + " "
		}

		str := fmt.Sprintf("%s\t%s\t%s", cmdata.Name, arguments, cmdata.Desc)
		fmt.Fprintln(w, str)
	}

	w.Flush()
}
