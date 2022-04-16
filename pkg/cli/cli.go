package cli

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
)

type CommandBase struct {
	programName string
	commandName string
}

type Arg struct {
	Name        string
	Description string
}

type Command interface {
	ProgName() string    // Name of program used for pretty-printing.
	CmdName() string     // Name of command.
	Overview() string    // Oneliner description of command.
	Arguments() *[]Arg   // Needed arguments to use command.
	Usage() string       // How to use the command.
	Description() string // Detailed description.
	Run([]string) error  // Run expects only args inteded for this command.
}

func checkCmdArguments(args []string, c Command) error {
	if len(args) == 0 {
		fmt.Println(GenerateUsage(c))
		return errors.New("zero arguments given")
	}

	if args[len(args)-1] == "--help" {
		fmt.Println(GenerateUsage(c))
		fmt.Print("Description:")
		fmt.Println(c.Description())
		return errors.New("help flag given")
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
