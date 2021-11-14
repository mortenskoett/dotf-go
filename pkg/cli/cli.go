package cli

import (
//	"fmt"
	"strings"
//	"text/tabwriter"

//	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type Command interface {
	Name() string
	Overview() string
	Arguments() map[string]string
	Usage() string
	Description() string

	// Run expects only args inteded for this command is given.
	Run([]string) error
}

func GenerateHelp(name, usage string, args map[string]string) string {
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
