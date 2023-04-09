/*
Package config contains functionality relating to the configuration file of dotf.
The parser is inspired by the TOML layout. See https://toml.io/en/ for format details.
Only a microscopic subset of the TOML v1.0.0 specification is implemented.
*/
package parsing

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mortenskoett/dotf-go/pkg/logging"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

// Defaults
var (
	userConfigDir, _ = os.UserConfigDir()
	defaultConfigDir = userConfigDir + "/dotf/config"
)

// Configurations that will be parsed from the config file
var (
	requiredConfigKeys = map[string]bool{
		"userspacedir": true,
		"dotfilesdir":  true,
		"syncdir":      true,
		"autosync":     false,
		"syncinterval": true,
	}
)

type DotfConfiguration struct {
	SyncDir          string // Git initialized directory that dotf should sync with remote
	UserspaceDir     string // Userspace dir is the root of the file hierachy dotf replicates
	DotfilesDir      string // Directory inside SyncDir containing same structure as userspace dir
	SyncIntervalSecs int    // Interval between syncing with remote using dotf-tray application
	AutoSync         bool   // If dotf-tray should autosync at given interval
}

/* Creates an empty Configuration with default values. */
func newConfiguration() *DotfConfiguration {
	return &DotfConfiguration{
		SyncDir:          "N/A",
		UserspaceDir:     "~/",
		DotfilesDir:      "N/A",
		SyncIntervalSecs: 120,
		AutoSync:         false,
	}
}

/*
* Config format:
* key0 = value0
* key1 = value1
* # is a comment
 */

// Parses the required dotf configuration file. Tries the given paths in order and otherwise
// fallbacks to default config location.
func ParseConfig(paths ...string) (*DotfConfiguration, error) {
	for _, p := range paths {
		if p == "" {
			continue
		}

		config, err := readConfigFrom(p)
		if err != nil {
			logging.Warn(fmt.Errorf("failed to parse config on path: %w", err))
			continue
		}
		return config, nil
	}
	defaultconfig, err := parseDefaultConfig()
	if err != nil {
		return nil, err
	}
	return defaultconfig, nil
}

func parseDefaultConfig() (*DotfConfiguration, error) {
	config, err := readConfigFrom(defaultConfigDir)
	if err == nil {
		return config, nil
	}
	return nil, &ParseConfigurationError{fmt.Sprintf("failed to parse config at default location: %s", defaultConfigDir)}
}

func readConfigFrom(path string) (*DotfConfiguration, error) {
	absPath, err := terminalio.GetAndValidateAbsolutePath(path)
	if err != nil {
		return nil, fmt.Errorf("path to config invalid: %w", err)
	}

	conf, err := readFromFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't load config: %w", err)
	}
	return conf, nil
}

// readFromFile parses and returns a representation of a config file found at 'path'.
func readFromFile(path string) (*DotfConfiguration, error) {
	config := newConfiguration()

	_, err := os.Stat(path)
	if err != nil {
		fmt.Println("configuration missing at", path)
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	keysToValues, err := parseTOMLFile(file)
	if err != nil {
		return nil, err
	}

	if err = checkRequiredKeys(keysToValues, requiredConfigKeys); err != nil {
		return config, err
	}

	err = buildConfiguration(config, keysToValues)
	if err != nil {
		return config, err
	}

	return config, nil
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

// Validate key values for required missing keys
func checkRequiredKeys(keysToValues map[string]string, requiredConfigKeys map[string]bool) error {
	for key, isRequired := range requiredConfigKeys {
		if isRequired {
			_, exists := keysToValues[key]
			if !exists {
				return &ConfigKeyNotFoundError{fmt.Sprint("missing key in configuration: ", key)}
			}
		}
	}
	return nil
}

/* parseTOMLFile parses the file found at 'file' and returns a key,value representation. */
func parseTOMLFile(file *os.File) (map[string]string, error) {
	parameterToValue := make(map[string]string)

	scanner := bufio.NewScanner(file)
	linenum := 0
	for scanner.Scan() {
		line := scanner.Text()
		nameAndValue := strings.SplitN(line, "=", 2)
		linenum++

		if strings.HasPrefix(nameAndValue[0], "#") {
			// Ignore outcommented lines.
			continue
		}

		// Didn't get both key and value
		if len(nameAndValue) < 2 {
			return nil, errors.New(
				fmt.Sprintf(
					"malformed key in configuration on line number: %d: %s",
					linenum,
					nameAndValue[0]))
		}

		parameter := strings.ToLower(strings.TrimFunc(nameAndValue[0], sanitize))
		value := strings.TrimFunc(nameAndValue[1], sanitize)
		parameterToValue[parameter] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return parameterToValue, nil
}

func sanitize(r rune) bool {
	return r == ' ' ||
		r == '\t' ||
		r == '\n' ||
		r == '\r' ||
		r == '"'
}

func buildConfiguration(config *DotfConfiguration, keyToValue map[string]string) error {
	// TODO: Could be made into json parsing in a smart go way
	for k, v := range keyToValue {
		switch k {
		case "dotfilesdir":
			config.DotfilesDir = expandTilde(v)
		case "userspacedir":
			config.UserspaceDir = expandTilde(v)
		case "syncdir":
			config.SyncDir = expandTilde(v)
		case "syncinterval":
			if v_num, err := strconv.Atoi(v); err != nil {
				return err
			} else {
				config.SyncIntervalSecs = v_num
			}
		case "autosync":
			config.AutoSync = true
		default:
			return &MalformedConfigurationError{fmt.Sprint(
				"malformed or unknown key encountered: ", k)}
		}
	}
	return nil
}

func expandTilde(path string) string {
	if strings.HasPrefix(path, "~/") {
		dirname, _ := os.UserHomeDir()
		path = filepath.Join(dirname, path[2:])
		return path
	}
	return path
}
