package argparse

import (
	"bytes"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mortenskoett/dotf-go/pkg/cli"
	"github.com/mortenskoett/dotf-go/pkg/logger"
)

func PrintBasicHelp(commands []cli.Command, programName, logo string) {
	printHeader(logo)
	printUsage(commands, programName)
}

func PrintFullHelp(commands []cli.Command, programName, logo string) {
	printHeader(logo)
	fmt.Println(`
Details:
	- Userspace describes where the symlinks are placed pointing into the dotfiles directory.
	- The dotfiles directory is where the actual configuration files are stored.
	- The folder structure in the dotfiles directory will match that of the userspace.`)
	printUsage(commands, programName)
}

func printHeader(logo string) {
	fmt.Println(logger.Color(logo, logger.Blue))
	fmt.Println("Dotfiles handler in Go.")
}

func printUsage(commands []cli.Command, programName string) {
	fmt.Println("\nUsage:", programName, "<command> <args> [--help]")
	fmt.Println("")

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 4, ' ', 0)

	fmt.Println("Commands:")

	// Print commands
	for _, cmd := range commands {
		buf := &bytes.Buffer{}
		for _, arg := range cmd.Arguments() {
			buf.WriteString("<")
			buf.WriteString(arg.Name)
			buf.WriteString(">")
			buf.WriteString("  ")
		}

		str := fmt.Sprintf("\t%s\t%s\t%s", cmd.CmdName(), buf.String(), cmd.Overview())
		fmt.Fprintln(w, str)
	}

	w.Flush()
	fmt.Println()
}
