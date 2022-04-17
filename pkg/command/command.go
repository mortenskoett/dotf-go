// Package command contains handling of all dotf operations given by cli arg
package command

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/mortenskoett/dotf-go/pkg/shared/global"
)

// ** ALL PROGRAM COMMANDS AVAILABLE BELOW ** //

// Contains the CLI Commands that are currently implemented in dotf. The commands are returned as
// functions so the name of the application can be given as param. The program name is used for
// pretty-printing.
var commands = map[string]Command{
	"add":  NewAddCommand(global.ProgramName, "add"),
	"move": NewMoveCommand(global.ProgramName, "move"),
}

type CommandBase struct {
	programName string
	commandName string
}

type Arg struct {
	Name        string
	Description string
}

type Command interface {
	ProgName() string             // Name of program used for pretty-printing.
	CmdName() string              // Name of command.
	Overview() string             // Oneliner description of the command.
	Arguments() *[]Arg            // Needed arguments to use the command.
	Usage() string                // How to use the command.
	Description() string          // Detailed description.
	Run(args *CliArguments) error // Attempt to runs the Command using the given args
}

// Contains CLI arguments parsed. Type placed here due to circular dependency.
type CliArguments struct {
	PosArgs []string // In order by input
	Flags   map[string]string
}

func NewCliArguments() *CliArguments {
	return &CliArguments{
		Flags: make(map[string]string),
	}
}

// Get copy of all Commands
func GetAllCommands() map[string]Command {
	return commands
}

// Creates a Command from command name and program name.
func ParseCommandName(cmdName string) (Command, error) {
	cmd, ok := commands[cmdName]
	if ok {
		return cmd, nil
	}
	return nil, fmt.Errorf("%s command does not exist.", cmdName)
}

func checkCmdArguments(args *CliArguments, c Command) error {
	if _, ok := args.Flags["help"]; ok {
		fmt.Println(GenerateUsage(c))
		fmt.Print("Description:")
		fmt.Println(c.Description())
		return &CmdErrorSuccess{"help flag given"}
	}

	if len(args.PosArgs) != len(*c.Arguments()) {
		fmt.Println(GenerateUsage(c))
		return fmt.Errorf(fmt.Sprintf(
			"%d arguments given, but %d required. Try adding --help.", len(args.PosArgs), len(*c.Arguments())))
	}

	return nil
}

// Generates a pretty-printed usage description of a Command
func GenerateUsage(c Command) string {
	var sb strings.Builder

	sb.WriteString("Name:\n\t")
	name := fmt.Sprintf("%s %s - %s", c.ProgName(), c.CmdName(), c.Overview())
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

// Displays a yes/no prompt to the user and returns the boolean value of the answer
func confirmByUser(question string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [Y(yes)/n(no)]\n", question)

		resp, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		resp = strings.TrimSpace(resp)

		if resp == "Y" || resp == "yes" {
			return true
		} else if resp == "n" || resp == "no" {
			return false
		}
	}
}
