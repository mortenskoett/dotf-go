package cli

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"bytes"

//	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type Arg struct {
	Name string
	Description string
}

type Command interface {
	Name() string	// Name of command.
	Overview() string // Oneliner description of command.
	Arguments() *[]Arg // Needed arguments to use command.
	Usage() string		// How to use the command.
	Description() string // Detailed description.
	Run([]string) error // Run expects only args inteded for this command.
}

func GenerateUsage(programName string, c Command) string {
	var sb strings.Builder

	sb.WriteString("Name:\n\t")
	name := fmt.Sprintf( "%s %s - %s", programName, c.Name(), c.Overview())
	sb.WriteString(name)

	sb.WriteString("\n\nUsage:\n\t")
	sb.WriteString(c.Usage())

	sb.WriteString("\n\nArguments:\n")

	// Print arguments.
	tabbuf := &bytes.Buffer{}
	w := new(tabwriter.Writer)
	w.Init(tabbuf, 0, 8, 8, ' ', 0)

	for _, arg := range *c.Arguments() {
		buf := &bytes.Buffer{}
		buf.WriteString("<")
		buf.WriteString(arg.Name)
		buf.WriteString(">")
		str := fmt.Sprintf("\t%s\t%s", buf, arg.Description)
		fmt.Fprintln(w, str)
	}

	w.Flush()
	sb.WriteString(tabbuf.String())

	return sb.String()
}

func GenerateHelp(c Command) string {
	var sb strings.Builder

	sb.WriteString("Name:")
	sb.WriteString("\n")

	sb.WriteString("Usage:")
	sb.WriteString("\n")

	sb.WriteString("Description:")
	sb.WriteString("\n")

	sb.WriteString("Arguments:")
	sb.WriteString("\n")
	return sb.String()
}
