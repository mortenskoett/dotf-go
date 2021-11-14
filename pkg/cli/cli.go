package cli

import (
//	"fmt"
	"strings"
//	"text/tabwriter"

//	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type CommandData struct {
	Name string
	Args map[string]string	// Taken arguments and their meaning.
	Desc string				// Short description of command.
	Logo string
}

type Command interface {
	Run([]string) error
	Data() CommandData
}

func GenerateHelp(cdata CommandData) string {
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

//	// Construct argument list
//	sb.WriteString(terminalio.Color(fmt.Sprintf("Usage: %s <to> <from>", cdata.Name), terminalio.Yellow))
//	sb.WriteString("\n")
//	sb.WriteString(terminalio.Color(fmt.Sprintf("Usage: %s help", cdata.Name), terminalio.Yellow))
//	sb.WriteString("\n\n")
//
//	tw := tabwriter.NewWriter(&sb, 0, 8, 4, ' ', 0)
//
//	fmt.Fprintln(tw, ("COMMAND\tARGS\tDESCRIPTION"))
//
//	// Formatting arguments.
//	var arguments string
//	for arg, _ := range cdata.Args {
//		arguments += "<" + arg + ">" + " "
//	}
//
//	str := fmt.Sprintf("%s\t%s\t%s", cdata.Name, arguments, cdata.Desc)
//
//	fmt.Fprint(tw, str)
//	tw.Flush()
