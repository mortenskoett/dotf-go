package argparse

import (
	"bytes"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mortenskoett/dotf-go/pkg/cli"
	"github.com/mortenskoett/dotf-go/pkg/logging"
)

func PrintBasicHelp(commands []cli.Command, programName, logo, version string) {
	printHeader(logo, version)
	printUsage(commands, programName)
}

func PrintFullHelp(commands []cli.Command, programName, logo, version string) {
	printHeader(logo, version)
	fmt.Println(`
Details:
	- Userspace describes where the symlinks are placed pointing into the dotfiles directory.
	- The dotfiles directory is where the actual configuration files are stored.
	- The folder structure in the dotfiles directory will match that of the userspace.`)
	printUsage(commands, programName)
}

func printHeader(logo, version string) {
	fmt.Println(logging.Color(logo, logging.Blue))
	fmt.Print("Dotfiles handler in Go.")
	fmt.Println(" Version:", version)
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
		if len(cmd.Arguments()) > 0 {
			for _, arg := range cmd.Arguments() {
				buf.WriteString("<")
				buf.WriteString(arg.Name)
				buf.WriteString(">")
				buf.WriteString("  ")
			}
		} else {
			buf.WriteString("-")
		}

		str := fmt.Sprintf("\t%s\t%s\t%s", cmd.CmdName(), buf.String(), cmd.Overview())
		fmt.Fprintln(w, str)
	}

	w.Flush()
	fmt.Println()
}
