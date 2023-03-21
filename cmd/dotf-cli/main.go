package main

import (
	"os"

	"github.com/mortenskoett/dotf-go/pkg/cli"
	"github.com/mortenskoett/dotf-go/pkg/logging"
	"github.com/mortenskoett/dotf-go/pkg/parsing"
)

const logo = `    _       _     __             _  _
 __| | ___ | |_  / _|  ___   __ | |(_)
/ _' |/ _ \|  _||  _| |___| / _|| || |
\__/_|\___/ \__||_|         \__||_||_|
`

// Inserted by build process using -ldflags
var programVersion = ""

// Available CLI commands
var (
	commands = []cli.Command{
		cli.NewAddCommand("add"),
		cli.NewInstallCommand("install"),
		cli.NewMigrateCommand("migrate"),
		cli.NewSyncCommand("sync"),
		cli.NewRevertCommand("revert"),
	}
)

func main() {
	run(os.Args)
}

func run(osargs []string) {
	// Create command env to manage commands
	executor := cli.NewCmdExecutor(commands)

	// Parse cli args
	cliargs, err := parsing.ParseCommandlineInput(os.Args)
	if err != nil {
		handleParsingError(err, executor.GetPrintableCmds())
	}

	// Parse dotf config
	config, err := parsing.ParseDotfConfig(cliargs.Flags)
	if err != nil {
		handleParsingError(err, executor.GetPrintableCmds())
	}

	// Create command
	cmd, err := executor.Load(cliargs, config)
	if err != nil {
		switch err.(type) {
		case *cli.CmdUnknownCommand:
			logging.Error(err)
		default:
			logging.Error("unknown executor load error:", err)
		}
		os.Exit(1)
	}

	// If no errors then run command
	err = cmd()
	if err != nil {
		switch err.(type) {
		case *cli.CmdHelpFlagError:
			logging.Ok(err)
		case *cli.CmdArgumentError:
			logging.Warn(err)
		case *cli.GitError:
			logging.Error(err)
		default:
			logging.Error("unknown command run error:", err)
		}
		os.Exit(1)
	}
	logging.Ok("All good. Bye!")
}

// Handles and terminates program according to error type
func handleParsingError(err error, cmds []cli.CommandPrintable) {
	if err != nil {
		switch err.(type) {
		case *parsing.ParseHelpFlagError:
			cli.PrintFullHelp(cmds, logo, programVersion)
			logging.Ok(err)
		case *parsing.ParseNoArgumentError:
			cli.PrintBasicHelp(cmds, logo, programVersion)
			logging.Ok(err)
		case *parsing.ParseInvalidArgumentError:
			logging.Warn(err)
		case *parsing.ParseConfigurationError:
			logging.Error(err)
		case *parsing.ParseError:
			logging.Error(err)
		default:
			logging.Error("unknown parser error:", err)
		}
		os.Exit(1)
	}
}
