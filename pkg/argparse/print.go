package argparse

import (
	"bytes"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mortenskoett/dotf-go/pkg/cli"
	"github.com/mortenskoett/dotf-go/pkg/constant"
	"github.com/mortenskoett/dotf-go/pkg/logger"
)

func printBasicHelp() {
	printHeader()
	printUsage()
}

func printFullHelp() {
	printHeader()
	fmt.Println(`
Details:
	- User space describes where the symlinks are placed pointing into the dotfiles directory.
	- The dotfiles directory is where the actual configuration files are stored.
	- The folder structure in the dotfiles directory will match that of the user space.`)
	printUsage()
}

func printHeader() {
	fmt.Println(logger.Color(constant.Logo, logger.Blue))
	fmt.Println("Dotfiles handler in Go.")
}

func printUsage() {
	fmt.Println("\nUsage:", constant.ProgramName, "<command> <args> [--help]")
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
