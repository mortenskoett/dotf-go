// The main dotf-go application entry point.
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/mortenskoett/dotf-go/pkg/cli"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

const (
	logo = `    _       _     __         __ _      
 __| | ___ | |_  / _|  ___  / _' | ___ 
/ _' |/ _ \|  _||  _| |___| \__. |/ _ \
\__/_|\___/ \__||_|         |___/ \___/
`
	programName string = "dotf-go"
)

var (
	// commands contains the CLI commands that are currently implemented in dotf.
	commands = map[string]cli.Command{
		"add":  cli.NewAddCommand(programName, "add"),
		"move": cli.NewMoveCommand(programName, "move"),
	}
)

func main() {
	log.SetFlags(0)
	log.SetPrefix(terminalio.Color("dotf-cli error: ", terminalio.Red))
	args := os.Args[1:]

	if len(args) > 0 {
		handleArguments(args)
	} else {
		printHeader()
		printUsage()
	}
}

func handleArguments(args []string) {
	input := args[0]
	count := len(args)

	if input == "" || input == "help" || input == "--help" || count == 0 {
		printHelp()
		return
	}

	cmd, err := parseCommand(input)
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Run(args[1:])
	if err != nil {
		log.Fatal(err)
	}
}

func parseCommand(input string) (cli.Command, error) {
	cmd, ok := commands[input]
	if ok {
		return cmd, nil
	}

	return nil, fmt.Errorf("%s command does not exist. Try adding --help.", input)
}

func printUsage() {
	fmt.Println("\nUsage: dotf-go <command> <args> [--help]")
	fmt.Println("")

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 4, ' ', 0)

	fmt.Println("Commands:")

	// Print commands
	for _, c := range commands {
		buf := &bytes.Buffer{}
		for _, arg := range *c.Arguments() {
			buf.WriteString("<")
			buf.WriteString(arg.Name)
			buf.WriteString(">")
			buf.WriteString("  ")
		}

		str := fmt.Sprintf("\t%s\t%s\t%s", c.CmdName(), buf.String(), c.Overview())
		fmt.Fprintln(w, str)
	}

	w.Flush()
}

func printHeader() {
	fmt.Println(terminalio.Color(logo, terminalio.Blue))
	fmt.Println("Dotfiles handler in Go.")
}

func printHelp() {
	printHeader()
	fmt.Println(`
Terminology:
	1) User space describes where the symlinks are placed pointing into the dotfiles directory.
	2) The dotfiles directory is where the actual configuration files are stored.
	3) The folder structure in the dotfiles directory will match that of the user space.`)
	printUsage()
}
