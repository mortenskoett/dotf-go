/*
Package config contains functionality relating to the configuration file of dotf.
The parser is inspired by the TOML layout. See https://toml.io/en/ for format details.
Only a microscopic subset of the TOML v1.0.0 specification is implemented.
*/
package parsing

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mortenskoett/dotf-go/pkg/logging"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

// Env vars
var (
	homedir, _   = os.UserHomeDir()
	configdir, _ = os.UserConfigDir()
	hostname, _  = os.Hostname()
)

// Defaults
var (
	defaultConfigDir   = configdir + "/dotf/config"
	defaultSyncDir     = homedir + "/dotfiles"
	defaultDotfilesDir = defaultSyncDir + "/distros/" + hostname
)

// Configurations that will be parsed from the config file
const (
	userspacedir     = "userspacedir"
	dotfilesdir      = "dotfilesdir"
	syncdir          = "syncdir"
	autosync         = "autosync"
	syncintervalsecs = "syncintervalsecs"
)

// Configurations that are required for dotf to function properly
var (
	requiredConfigKeys = map[string]bool{
		userspacedir:     true,
		dotfilesdir:      true,
		syncdir:          true,
		autosync:         false,
		syncintervalsecs: true,
	}
)

type ConfigMetadata struct {
	Filepath string `json:"filepath"` // Not configurable
}

type DotfConfiguration struct {
	*ConfigMetadata
	UserspaceDir     string `json:"userspacedir"`     // Userspace dir is the root of the file hierachy dotf replicates
	DotfilesDir      string `json:"dotfilesdir"`      // Directory inside SyncDir containing same structure as userspace dir
	SyncDir          string `json:"syncdir"`          // Git initialized directory that dotf should sync with remote
	AutoSync         bool   `json:"autosync"`         // If dotf-tray should autosync at given interval
	SyncIntervalSecs int    `json:"syncintervalsecs"` // Interval between syncing with remote using dotf-tray application
}

/* Creates a basic sensible Configuration with default values. */
func NewSensibleConfiguration() *DotfConfiguration {
	return &DotfConfiguration{
		ConfigMetadata: &ConfigMetadata{
			Filepath: defaultConfigDir,
		},
		UserspaceDir:     homedir,
		DotfilesDir:      defaultDotfilesDir,
		SyncDir:          defaultSyncDir,
		AutoSync:         false,
		SyncIntervalSecs: 3600,
	}
}

func NewEmptyConfiguration() *DotfConfiguration {
	return &DotfConfiguration{
		ConfigMetadata: &ConfigMetadata{
			Filepath: "",
		},
		UserspaceDir:     "",
		DotfilesDir:      "",
		SyncDir:          "",
		AutoSync:         false,
		SyncIntervalSecs: 3600,
	}
}

// Convert a configuration to a map, used for easier serialization of configuration
func ConvertConfigToMap(conf *DotfConfiguration) (map[string]string, error) {
	// from conf -> json
	b, err := json.Marshal(conf)
	if err != nil {
		return nil, err
	}

	// from json -> map
	confmap := make(map[string]any)
	err = json.Unmarshal(b, &confmap)
	if err != nil {
		return nil, err
	}

	// use only required keys
	for k := range confmap {
		if _, ok := requiredConfigKeys[k]; !ok {
			delete(confmap, k)
		}
	}

	// validate
	for k := range requiredConfigKeys {
		if _, ok := confmap[k]; !ok {
			return nil, fmt.Errorf("failed to validate keys when creating config map: %v", err)
		}
	}

	// conv to string map
	strmap := make(map[string]string)
	for k, v := range confmap {
		switch t := v.(type) {
		case bool:
			s := fmt.Sprintf("%t", t)
			strmap[k] = s
		case float64:
			s := fmt.Sprintf("%d", int(t))
			strmap[k] = s
		default:
			strmap[k] = v.(string)
		}
	}
	return strmap, nil
}

// Creates a slice of bytes that can be serialized to a file and used as a valid config
func CreateSerializableConfig(keyvals map[string]string) []byte {
	var builder strings.Builder
	for k, v := range keyvals {
		builder.WriteString(k)
		builder.WriteString(" = ")
		builder.WriteString(v)
		builder.WriteString("\n")
	}
	return []byte(builder.String())
}

/*
* Config format:
* key0 = value0
* key1 = value1
* # is a comment
 */

// Parses the required dotf configuration file. Tries the given paths in order or otherwise
// fallbacks to default config location.
func ParseConfig(paths ...string) (*DotfConfiguration, error) {
	for _, p := range paths {
		if p == "" {
			continue
		}

		config, err := readConfig(p)
		if err != nil {
			logging.Warn(fmt.Errorf("failed to parse config on path: %w", err))
			continue
		}
		return config, nil
	}

	config, err := readConfig(defaultConfigDir)
	if err != nil {
		return NewEmptyConfiguration(), &ParseDefaultConfigurationError{fmt.Sprintf("%v", err)}
	}
	return config, nil
}

func readConfig(path string) (*DotfConfiguration, error) {
	absPath, err := terminalio.GetAndValidateAbsolutePath(path)
	if err != nil {
		return nil, fmt.Errorf("path to config invalid: %w", err)
	}

	conf, err := parseConfig(absPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't load config: %w", err)
	}
	return conf, nil
}

// parseConfig parses and returns a representation of a config file found at 'path'.
func parseConfig(path string) (*DotfConfiguration, error) {
	config := NewEmptyConfiguration()

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

	if err = validateKeys(keysToValues, requiredConfigKeys); err != nil {
		return config, err
	}

	err = buildConfiguration(config, keysToValues)
	if err != nil {
		return config, err
	}

	config.Filepath = path
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

// Validate key values for required but potentially missing keys
func validateKeys(keysToValues map[string]string, requiredConfigKeys map[string]bool) error {
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
		case dotfilesdir:
			config.DotfilesDir = expandTilde(v)
		case userspacedir:
			config.UserspaceDir = expandTilde(v)
		case syncdir:
			config.SyncDir = expandTilde(v)
		case syncintervalsecs:
			if v_num, err := strconv.Atoi(v); err != nil {
				return err
			} else {
				config.SyncIntervalSecs = v_num
			}
		case autosync:
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
