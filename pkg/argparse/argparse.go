package argparse

import (
	"fmt"
	"os"
	"strings"

	"github.com/mortenskoett/dotf-go/pkg/cli"
	"github.com/mortenskoett/dotf-go/pkg/config"
	"github.com/mortenskoett/dotf-go/pkg/logger"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
	"github.com/mortenskoett/dotf-go/pkg/utils"
)

// Flags required to contain a value like 'exec cmd --flag value'. This is maintained by the parsing
// routine which will fail.
type ValueFlags []string

var valueflags ValueFlags = []string{"config"}

// Parses the CLI input arguments and the dotf configuration. Expects complete input argument line.
func Parse(osargs []string) (*cli.CliArguments, *config.DotfConfiguration, error) {
	cliargs, err := parseCliArguments(osargs)
	if err != nil {
		return nil, nil, err
	}

	conf, err := parseDotfConfig(cliargs.Flags)
	if err != nil {
		return nil, nil, err
	}

	return cliargs, conf, nil
}

// Parses CLI arguments into positional args and flags
func parseCliArguments(osargs []string) (*cli.CliArguments, error) {
	// Ignore executable name
	args := osargs[1:]

	if len(args) < 1 {
		return nil, &ParseNoArgumentError{"no arguments given"}
	}

	cmdName := args[0]
	count := len(args)
	if cmdName == "" || cmdName == "-h" || cmdName == "--h" || cmdName == "help" || cmdName == "--help" || count == 0 {
		return nil, &ParseHelpFlagError{"showing full help."}
	}

	cliargs, err := parseArgsAndFlags(args)
	if err != nil {
		return nil, &ParseError{fmt.Sprintf("failed to parse input: %s", err)}
	}

	return cliargs, nil
}

// Parses cli command and arguments without judgement about whether arguments are fit for Command.
func parseArgsAndFlags(osargs []string) (*cli.CliArguments, error) {
	cliarg := cli.NewCliArguments()

	cmdName := osargs[0]
	args := osargs[1:]

	cliarg.CmdName = cmdName

	parsePositionalInto(args, cliarg)

	if err := parseFlagsInto(args, valueflags, cliarg); err != nil {
		return nil, err
	}

	return cliarg, nil
}

// Parses only positional args and stops at the first flag e.g. '--flag'. The args are added to the
// supplied cli.Arguments.
func parsePositionalInto(args []string, cliarg *cli.CliArguments) {
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			break
		} else {
			cliarg.PosArgs = append(cliarg.PosArgs, arg)
		}
	}
}

// Parses only flags but both boolean and value holding flags The flags are added to the supplied
// cli.Arguments.
func parseFlagsInto(args []string, valueflags ValueFlags, cliarg *cli.CliArguments) error {
	var currentflag string

	for i, arg := range args {
		// previous arg was a value containing flag
		if currentflag != "" {
			if strings.HasPrefix(arg, "--") {
				// next is also flag
				return &ParseError{fmt.Sprintf(
					"given flag '--%s' must be followed by a value not a flag", currentflag)}
			}

			cliarg.Flags[currentflag] = arg
			currentflag = ""
			continue
		}

		// flags
		if strings.HasPrefix(arg, "--") {
			if i == len(args)-1 {
				// if last element
				return &ParseError{fmt.Sprintf(
					"given flag '%s' must be followed by a value, but was empty", arg)}
			}

			flag := strings.ReplaceAll(arg, "--", "")

			if utils.Contains(valueflags, flag) {
				// with value
				currentflag = flag

			} else {
				// no value
				cliarg.Flags[flag] = flag
			}

		}
	}
	return nil
}

// Parses the required dotf configuration file.
// 1. First --config <path> flag is tried and used in case it is valid
// 2. Then ${HOME}/.config/dotf/config is tried
// 3. If both fails a specifc parse config error is returned
func parseDotfConfig(flags map[string]string) (*config.DotfConfiguration, error) {
	if path, ok := flags["config"]; ok {
		config, err := readConfigFrom(path)
		if err == nil {
			logger.LogSuccess("Successfully parsed config at:", path)
			return config, nil
		}
		logger.LogWarn(fmt.Errorf("failed to parse config path from flag: %w", err))
	}

	configPath, _ := os.UserConfigDir()
	defaultPath := configPath + "/dotf/config"

	config, err := readConfigFrom(defaultPath)
	if err == nil {
		logger.LogSuccess("Successfully parsed config at default path: ", defaultPath)
		return config, nil
	}
	logger.LogWarn(fmt.Errorf("failed to parse config at default location: %w", err))

	return nil, &ParseConfigurationError{"no valid dotf configuration found."}
}

func readConfigFrom(path string) (*config.DotfConfiguration, error) {
	absPath, err := terminalio.GetAndValidateAbsolutePath(path)
	if err != nil {
		return nil, fmt.Errorf("path to config invalid: %w", err)
	}

	conf, err := config.ReadFromFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't load config: %w", err)
	}
	return &conf, nil
}
