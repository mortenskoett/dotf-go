// Package command contains handling of all dotf operations given by cli arg
package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/mortenskoett/dotf-go/pkg/config"
	"github.com/mortenskoett/dotf-go/pkg/logging"
)

// ** ALL PROGRAM COMMANDS AVAILABLE BELOW ** //

// CommandFunc defines a function that given the name of the executable (most likely dotf-go) will
// return a valid Command.
type CommandFunc = func(execName string) Command

// Contains the CLI Commands that are currently implemented in dotf. The commands are returned as
// functions so the name of the application can be given as param. The program name is used for
// pretty-printing.
var commands = map[string]CommandFunc{
	"add":     func(pname string) Command { return NewAddCommand(pname, "add") },
	"install": func(pname string) Command { return NewInstallCommand(pname, "install") },
	"move":    func(pname string) Command { return NewMoveCommand(pname, "move") },
	"sync":    func(pname string) Command { return NewSyncCommand(pname, "sync") },
	"revert":  func(pname string) Command { return NewRevertCommand(pname, "revert") },
}

// Contains basic program info for each Command
type CommandBase struct {
	programName string
	commandName string
}

// Defines a required argument for a specific Command
type Arg struct {
	Name        string
	Description string
}

// Command is a definition of a main operation taking a number of cli args to work on
type Command interface {
	ProgName() string    // Name of program used for pretty-printing.
	CmdName() string     // Name of command.
	Overview() string    // Oneliner description of the command.
	Arguments() []Arg    // Needed arguments to use the command.
	Usage() string       // How to use the command.
	Description() string // Detailed description.
	Run(args *CliArguments,
		conf *config.DotfConfiguration) error // Attempt to run the Command using the given args and config
}

// Parsed CLI arguments
type CliArguments struct {
	CmdName string            // The first arg read after executable name
	PosArgs []string          // Positional args read after command in order
	Flags   map[string]string // Flags given after positional args in input
}

func NewCliArguments() *CliArguments {
	return &CliArguments{
		Flags: make(map[string]string),
	}
}

// Get copy of all available Commands. Obs: Ineffective implementation.
func GetAvailableCommands(programName string) []Command {
	cmds := make([]Command, 0, len(commands))
	for _, cmdf := range commands {
		cmds = append(cmds, cmdf(programName))
	}
	sort.SliceStable(cmds, func(i, j int) bool {
		return cmds[i].CmdName() < cmds[j].CmdName()
	})
	return cmds
}

// Creates a Command or errors
func CreateCommand(programName, cmdName string) (Command, error) {
	cmdfunc, err := parseToCommandFunc(cmdName)
	if err != nil {
		return nil, &CmdUnknownCommand{fmt.Sprintf("try --help for available commands: %s", err)}
	}
	return cmdfunc(programName), nil
}

// Parses a Command name to a CommandFunc or errors
func parseToCommandFunc(cmdName string) (CommandFunc, error) {
	cmdfunc, ok := commands[cmdName]
	if ok {
		return cmdfunc, nil
	}
	return nil, &CmdArgumentError{fmt.Sprintf("%s command does not exist.", cmdName)}
}

// Validates and handles the given Arguments generally against the Command and errors if not valid
func validateCliArguments(args *CliArguments, c Command) error {

	if _, ok := args.Flags["help"]; ok {
		fmt.Println(GenerateUsage(c))
		fmt.Print("Description:")
		fmt.Println(c.Description())
		return &CmdHelpFlagError{"help flag given"}
	}

	if len(args.PosArgs) != len(c.Arguments()) {
		fmt.Println(GenerateUsage(c))
		return &CmdArgumentError{fmt.Sprintf(
			"%d arguments given, but %d required. Try adding --help.", len(args.PosArgs), len(c.Arguments()))}
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

	for _, arg := range c.Arguments() {
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
