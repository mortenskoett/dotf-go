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

func printBasicHelp(execName string) {
	printHeader(execName)
	printUsage(execName)
}

func printFullHelp(execName string) {
	printHeader(execName)
	fmt.Println(`
Details:
	- User space describes where the symlinks are placed pointing into the dotfiles directory.
	- The dotfiles directory is where the actual configuration files are stored.
	- The folder structure in the dotfiles directory will match that of the user space.`)
	printUsage(execName)
}

func printHeader(execName string) {
	fmt.Println(logger.Color(constant.Logo, logger.Blue))
	fmt.Println("Dotfiles handler in Go.")
}

func printUsage(execName string) {
	fmt.Println("\nUsage:", execName, "<command> <args> [--help]")
	fmt.Println("")

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 4, ' ', 0)

	fmt.Println("Commands:")

	// Print commands
	for _, cmdFunc := range cli.GetCommandFuncs() {
		c := cmdFunc(execName)
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
