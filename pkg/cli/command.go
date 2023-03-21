// Package command contains handling of all dotf operations given by cli arg
package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/mortenskoett/dotf-go/pkg/logging"
	"github.com/mortenskoett/dotf-go/pkg/parsing"
)

const (
	programName string = "dotf-cli"
)

// Defines an argument for a specific Command
type arg struct {
	name        string
	description string
}

type cmd struct {
	name        string // Name of command.
	overview    string // Oneliner description of the command.
	description string // Detailed description.
	usage       string // How to use the command.
	arguments   []arg  // Required arguments in order to use the command.
}

// CommandPrintable is an interface used where it is necessary to print the command details
type CommandPrintable interface {
	CmdName() string     // Name of command.
	Overview() string    // Oneliner description of the command.
	Description() string // Detailed description.
	Usage() string       // How to use the command.
	Arguments() []arg    // Required arguments in order to use the command.
}

// CommandRunner is a definition of a main operation taking a number of cli args to work on
type CommandRunner interface {
	Run(args *parsing.CommandLineInput,
		conf *parsing.DotfConfiguration) error // Run the Command using the given args and config
}

type Command interface {
	CommandPrintable
	CommandRunner
}

// Validates and handles the given Arguments generally against the Command and errors if not valid
// TODO: This should be rewritten
func validateCliArguments(args *parsing.CommandLineInput, c CommandPrintable) error {
	if _, ok := args.Flags.BoolFlags["help"]; ok {
		fmt.Println(generateUsage(c))
		fmt.Print("Description:")
		fmt.Println(c.Description())
		return &CmdHelpFlagError{"help flag given"}
	}

	if len(args.PositionalArgs) != len(c.Arguments()) {
		fmt.Println(generateUsage(c))
		return &CmdArgumentError{fmt.Sprintf(
			"%d arguments given, but %d required. Try adding --help.", len(args.PositionalArgs), len(c.Arguments()))}
	}
	return nil
}

// Generates a pretty-printed usage description of a Command
func generateUsage(c CommandPrintable) string {
	var sb strings.Builder

	sb.WriteString("Name:\n\t")
	name := fmt.Sprintf("%s %s - %s", programName, c.CmdName(), c.Overview())
	sb.WriteString(name)

	sb.WriteString("\n\nUsage:\n\t")
	sb.WriteString(c.Usage())

	sb.WriteString("\n\nArguments:\n")

	// Print arguments.
	tabbuf := &bytes.Buffer{}
	w := new(tabwriter.Writer)
	w.Init(tabbuf, 0, 8, 8, ' ', 0)

	for _, arg := range c.Arguments() {
		buf := &bytes.Buffer{}
		buf.WriteString("<")
		buf.WriteString(arg.name)
		buf.WriteString(">")
		str := fmt.Sprintf("\t%s\t%s", buf, arg.description)
		fmt.Fprintln(w, str)
	}

	w.Flush()
	sb.WriteString(tabbuf.String())

	return sb.String()
}

// Displays a yes/no prompt to the user and returns the boolean value of the answer
func confirmByUser(question string) bool {
	reader := bufio.NewReader(os.Stdin)

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
