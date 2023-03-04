package argparse

import (
	"fmt"
	"os"
	"strings"

	"github.com/mortenskoett/dotf-go/pkg/cli"
	"github.com/mortenskoett/dotf-go/pkg/config"
	"github.com/mortenskoett/dotf-go/pkg/logging"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

var (
	userConfigDir, _ = os.UserConfigDir()
	defaultConfigDir = userConfigDir + "/dotf/config"
)

// Type alias used to express flags by the parser
type Flags = map[string]string

// Specific flags required to contain a value like 'exec cmd --flag value'. This is maintained by
// the parsing routine which will fail.
type ValueFlags []string

var vflags ValueFlags = []string{"config"}

// Parses the CLI input arguments and the Dotf configuration and returns potential errors
func Parse(osargs []string) (*cli.CliArguments, *config.DotfConfiguration, error) {
	cliargs, clierr := parseCliArguments(osargs, vflags)

	var conferr error
	var conf *config.DotfConfiguration
	if cliargs != nil {
		conf, conferr = ParseDotfConfig(cliargs.Flags)
	} else {
		conf, conferr = ParseDotfConfig(nil)
	}

	// Fail on configuration error first
	if conferr != nil {
		return nil, nil, conferr
	}

	if clierr != nil {
		return nil, nil, clierr
	}

	return cliargs, conf, nil
}

// Parses CLI arguments into positional args and flags
func parseCliArguments(osargs []string, vflags ValueFlags) (*cli.CliArguments, error) {
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

	cliargs, err := parseArgsAndFlags(args, vflags)
	if err != nil {
		return nil, &ParseError{fmt.Sprintf("failed to parse input: %s", err)}
	}

	return cliargs, nil
}

// Parses cli command and arguments without judgement about whether arguments are fit for Command.
func parseArgsAndFlags(osargs []string, vflags ValueFlags) (*cli.CliArguments, error) {
	cliarg := cli.NewCliArguments()

	cmdName := osargs[0]
	args := osargs[1:]

	cliarg.CmdName = cmdName
	cliarg.PosArgs = parsePositionalArgs(args)

	flags, err := ParseFlags(args, vflags)
	if err != nil {
		return nil, err
	}

	cliarg.Flags = flags

	return cliarg, nil
}

// Parses only positional args and stops at the first flag e.g. '--flag'. The args are added to the
// supplied cli.Arguments.
func parsePositionalArgs(args []string) (posArgs []string) {
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			break
		} else {
			posArgs = append(posArgs, arg)
		}
	}
	return
}

// Parses only flags but both boolean and value holding flags The flags are added to the supplied
// cli.Arguments. E.g. --config <path> and --help
func ParseFlags(args []string, valueflags ValueFlags) (flags Flags, err error) {
	flags = Flags{}

	var currentflag string
	for i, arg := range args {
		// previous arg was a value containing flag
		if currentflag != "" {
			if strings.HasPrefix(arg, "--") {
				// next is also flag
				return flags, &ParseError{fmt.Sprintf(
					"given flag '--%s' must be followed by a value not a flag", currentflag)}
			}

			flags[currentflag] = arg
			currentflag = ""
			continue
		}

		// flags
		if strings.HasPrefix(arg, "--") {
			flag := strings.ReplaceAll(arg, "--", "")
			isValueFlag := containsString(valueflags, flag)

			if i == len(args)-1 && isValueFlag {
				// if last element
				return flags, &ParseError{fmt.Sprintf(
					"given flag '%s' must be followed by a value, but was empty", arg)}

			} else if isValueFlag {
				// with value
				currentflag = flag

			} else {
				// no value
				flags[flag] = flag
			}

		}
	}
	return
}

// Parses the required dotf configuration file.
// 1. If flags not nil then --config <path> flag is tried and used in case it is valid
// 2. Then ${HOME}/.config/dotf/config is tried
// 3. If both fails a specifc parse config error is returned
func ParseDotfConfig(flags map[string]string) (*config.DotfConfiguration, error) {
	if flags != nil { // Only try config pointed to by flags if any flags
		if path, ok := flags["config"]; ok {
			config, err := readConfigFrom(path)
			if err == nil {
				return config, nil
			}
			logging.Warn(fmt.Errorf("failed to parse config path from flag: %w", err))
		}
	}

	config, err := readConfigFrom(defaultConfigDir)
	if err == nil {
		return config, nil
	}
	logging.Error(fmt.Errorf("failed to parse config at default location: %w", err))

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
	return conf, nil
}

// Returns true if sl contains str.
func containsString(sl []string, str string) bool {
	for _, e := range sl {
		if str == e {
			return true
		}
	}
	return false
}
