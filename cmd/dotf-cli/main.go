// The main dotf-go application entry point.
package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"bytes"

	"github.com/mortenskoett/dotf-go/pkg/terminalio"
	"github.com/mortenskoett/dotf-go/pkg/cli"
)

const (
	logo =
`    _       _     __                             
 __| | ___ | |_  / _|
/ _  |/ _ \|  _||  _|
\__/_|\___/ \__||_|  
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
	count := len(args)

	if input == "help" {
		printHelp()
		return
	}

	cmd, err := getCommand(input)
		if err != nil {
			log.Fatal(err)
	}

	if args[count-1] == "--help" {
		cmd.Help(name of program, )		// TODO: Start here
	}

	err = cmd.Run(args[1:])
	if err != nil {
		log.Fatal(err)
	}
}

func printHelp() {
	fmt.Println(terminalio.Color(logo, terminalio.Blue))
	fmt.Println("Usage: dotf-go <command> <args> [--help]")
	fmt.Println("")

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 4, ' ', 0)

	fmt.Println("Commands:")

	// Print commands.
	for _, c := range commands {
		cdata := c.Data()

		buf := &bytes.Buffer{}
		for arg, _ := range cdata.Args {
			buf.WriteString(arg)
			buf.WriteByte(' ')
		}

		str := fmt.Sprintf("\t%s\t%s\t%s", cdata.Name, buf.String(), cdata.Desc)
		fmt.Fprintln(w, str)
	}

	w.Flush()
}
