/*
Package config contains functionality relating to the configuration file of dotf.
The parser is inspired by the TOML layout. See https://toml.io/en/ for format details.
Only a microscopic subset of the TOML v1.0.0 specification is implemented.
*/
package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type DotfConfiguration struct {
	SyncDir           string // Git initialized directory that dotf should sync with remote
	UserspaceDir      string // Userspace dir is the root of the file hierachy dotf replicates
	DotfilesDir       string // Directory inside SyncDir containing same structure as userspace dir
	UpdateIntervalSec int    // Interval between syncing with remote using dotf-tray application
}

/* Creates an empty Configuration with default values. */
func NewConfiguration() DotfConfiguration {

	return DotfConfiguration{
		SyncDir:           "N/A",
		UserspaceDir:      "~/",
		DotfilesDir:       "N/A",
		UpdateIntervalSec: 120,
	}
}

/*
* Config format:
* key0 = value0
* key1 = value1
* # is a comment
 */

/* ReadFromFile parses and returns a representation of a config file found at 'absPath'. */
func ReadFromFile(path string) (DotfConfiguration, error) {
	config := NewConfiguration()

	_, err := os.Stat(path)
	if err != nil {
		fmt.Println("configuration missing at", path)
		return config, err
	}

	file, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer file.Close()

	parametersToValues, err := parseTOMLFile(file)
	if err != nil {
		return config, err
	}

	err = buildConfiguration(&config, &parametersToValues)
	if err != nil {
		return config, err
	}

	return config, nil
}

/* parseTOMLFile parses the file found at 'file' and returns a key,value representation. */
func parseTOMLFile(file *os.File) (map[string]string, error) {
	parameterToValue := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		nameAndValue := strings.SplitN(line, "=", 2)

		if strings.HasPrefix(nameAndValue[0], "#") {
			// Ignore outcommented lines.
			continue
		}

		if len(nameAndValue) < 2 {
			return nil, errors.New("malformed parameter line in configuration on line: " + nameAndValue[0])
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

func buildConfiguration(config *DotfConfiguration, keyToValue *map[string]string) error {
	for k, v := range *keyToValue {
		switch k {
		case "dotfilesdir":
			config.DotfilesDir = expandTilde(v)
		case "userspacedir":
			config.UserspaceDir = expandTilde(v)
		case "syncdir":
			config.SyncDir = expandTilde(v)
		case "updateintervalsec":
			if v_num, err := strconv.Atoi(v); err != nil {
				return err
			} else {
				config.UpdateIntervalSec = v_num
			}
		default:
			return fmt.Errorf("malformed or unknown key encountered in configuration for: %s", k)
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
