package cli

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type CommandData struct {
	Name string
	Args map[string]string	// Taken arguments and their meaning.
	Desc string			// Short description of command.
}

type Command interface {
	Run([]string) error
	Data() CommandData
}

func BuildUsageText(cdata CommandData) string {
	var sb strings.Builder
	w := tabwriter.NewWriter(&sb, 0, 8, 4, ' ', 0)

	// TODO: Fix these arguments to be used from input parameter
	// TODO: Implement the individual design as seen in sketch file
	sb.WriteString(terminalio.Color(fmt.Sprintf("Usage: %s <to> <from>", cdata.Name), terminalio.Yellow))
	sb.WriteString("\n")
	sb.WriteString(terminalio.Color(fmt.Sprintf("Usage: %s help", cdata.Name), terminalio.Yellow))
	sb.WriteString("\n")
	sb.WriteString("\n")

	fmt.Fprintln(w, ("COMMAND\tARGS\tDESCRIPTION"))

	// Formatting arguments.
	var arguments string
	for arg, _ := range cdata.Args {
		arguments += "<" + arg + ">" + " "
	}

	str := fmt.Sprintf("%s\t%s\t%s", cdata.Name, arguments, cdata.Desc)

	fmt.Fprint(w, str)
	w.Flush()
	return sb.String()
}

