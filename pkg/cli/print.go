package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/mortenskoett/dotf-go/pkg/logging"
)

const (
	// Program name used for printing
	programName string = "dotf-cli"
)

func PrintBasicHelp[T CommandPrintable](commands []T, logo, version string) {
	printHeader(logo, version)
	printUsage(commands, programName)
}

// Encapsulates basic help with a more full context
func PrintFullHelp[T CommandPrintable](commands []T, logo, version string) {
	printHeader(logo, version)
	fmt.Println(`
Details:
  - Userspace describes where the symlinks are placed pointing into the dotfiles directory.
  - The dotfiles directory is where the actual configuration files are stored.
  - The folder structure in the dotfiles directory will match that of the userspace.`)
	printUsage(commands, programName)
}

func PrintCommandHelp(c CommandPrintable) {
	fmt.Println(generateUsage(c))
	fmt.Print("Description:")
	fmt.Println(c.getDescription())
}

// Prints program header
func printHeader(logo, version string) {
	fmt.Println(logging.Color(logo, logging.Blue))
	fmt.Print("Dotfiles handler in Go.")
	fmt.Println(" Version:", version)
}

func printUsage[T CommandPrintable](commands []T, programName string) {
	fmt.Println("\nUsage:", programName, "<command> <args> [--help]")
	fmt.Println("")

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 4, ' ', 0)

	fmt.Println("Commands:")

	// Print commands
	for _, c := range commands {
		buf := &bytes.Buffer{}
		if len(c.getArgs()) > 0 {
			for _, arg := range c.getArgs() {
				buf.WriteString("<")
				buf.WriteString(arg.Name)
				buf.WriteString(">")
				buf.WriteString("  ")
			}
		} else {
			buf.WriteString("-")
		}

		str := fmt.Sprintf("\t%s\t%s\t%s", c.getName(), buf.String(), c.getOverview())
		fmt.Fprintln(w, str)
	}

	w.Flush()
	fmt.Println()
}

// Generates a pretty-printed usage description of a Command
func generateUsage(c CommandPrintable) string {
	var sb strings.Builder

	sb.WriteString("Name:\n\t")
	name := fmt.Sprintf("%s %s - %s", programName, c.getName(), c.getOverview())
	sb.WriteString(name)

	sb.WriteString("\n\nUsage:\n\t")
	sb.WriteString(c.getUsage())

	// TODO: Fix duplicated code for printing of args and flags below.

	// Print arguments.
	sb.WriteString("\n\nArguments:\n")
	tabbuf := &bytes.Buffer{}
	w := new(tabwriter.Writer)
	w.Init(tabbuf, 0, 8, 8, ' ', 0)

	for _, arg := range c.getArgs() {
		buf := &bytes.Buffer{}
		buf.WriteString("<")
		buf.WriteString(arg.Name)
		buf.WriteString(">")
		str := fmt.Sprintf("\t%s\t%s", buf, arg.Description)
		fmt.Fprintln(w, str)
	}
	w.Flush()
	sb.WriteString(tabbuf.String())

	// Print flags.
	if len(c.getFlags()) > 0 {
		sb.WriteString("\nFlags:\n")
		tabbuf = &bytes.Buffer{}
		w = new(tabwriter.Writer)
		w.Init(tabbuf, 0, 8, 8, ' ', 0)

		for _, flag := range c.getFlags() {
			buf := &bytes.Buffer{}
			buf.WriteString("<")
			buf.WriteString(flag.Name)
			buf.WriteString(">")
			str := fmt.Sprintf("\t%s\t%s", buf, flag.Description)
			fmt.Fprintln(w, str)
		}
		w.Flush()
		sb.WriteString(tabbuf.String())
	}

	return sb.String()
}

type UserInteractor interface {
	ConfirmByUser(question string) bool
}

type StdInUserInteractor struct {
}

// Implements UserInteractor.
func (s StdInUserInteractor) ConfirmByUser(question string) bool {
	return ConfirmByUser(question, os.Stdin)
}

// Displays a yes/no prompt to the user and returns the boolean value of the answer. Stdin is
// parameterized to make the function testable.
func ConfirmByUser(question string, stdin io.Reader) bool {
	reader := bufio.NewReader(stdin)

	for {
		logging.Warn(question)
		logging.Input("[Y(yes)/n(no)]")

		resp, err := reader.ReadString('\n')
		if err != nil {
			logging.Fatal(err)
		}

		resp = strings.TrimSpace(resp)

		if resp == "Y" || resp == "yes" {
			return true
		} else if resp == "n" || resp == "no" {
			return false
		}
	}
}
