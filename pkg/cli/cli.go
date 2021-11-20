package cli

import (
//	"fmt"
	"strings"
//	"text/tabwriter"

//	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type Arg struct {
	Name string
	Description string
}

type Command interface {
	Name() string
	Overview() string
	Arguments() *[]Arg
	Usage() string
	Description() string

	// Run expects only args inteded for this command.
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
