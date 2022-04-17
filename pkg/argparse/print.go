package argparse

import (
	"bytes"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mortenskoett/dotf-go/pkg/cli"
	"github.com/mortenskoett/dotf-go/pkg/shared/global"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

func printBasicHelp() {
	printHeader()
	printUsage()
}

func printFullHelp() {
	printHeader()
	fmt.Println(`
Terminology:
	1) User space describes where the symlinks are placed pointing into the dotfiles directory.
	2) The dotfiles directory is where the actual configuration files are stored.
	3) The folder structure in the dotfiles directory will match that of the user space.`)
	printUsage()
}

func printHeader() {
	fmt.Println(terminalio.Color(global.Logo, terminalio.Blue))
	fmt.Println("Dotfiles handler in Go.")
}

func printUsage() {
	fmt.Println("\nUsage:", global.ProgramName, "<command> <args> [--help]")
	fmt.Println("")

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 4, ' ', 0)

	fmt.Println("Commands:")

	// Print commands
	for _, c := range cli.GetAllCommands() {
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
	fmt.Println()
}
