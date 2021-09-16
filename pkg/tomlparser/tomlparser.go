/*
Very basic TOML parser. See https://toml.io/en/ for format details.
Only a microscopic subset of the TOML v1.0.0 specification is implemented.
*/
package tomlparser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Configuration struct {
	RemoteURL         string
	DotFilesDir       string
	UpdateIntervalSec int
}

/* Creates an empty Configuration with default values. */
func NewConfiguration() Configuration {
	return Configuration{
		"https://www.github.de/someone/doesntexist",
		"$HOME/dotfiles/",
		10,
	}
}

/* ReadConfigurationFile parses and returns a representation of a *.toml file found at 'absPath'. */
func ReadConfigurationFile(absPath string) (Configuration, error) {
	config := NewConfiguration()

	_, err := os.Stat(absPath)
	if err != nil {
		fmt.Println("config.toml missing at", absPath)
		return config, err
	}

	file, err := os.Open(absPath)
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

		parameter := strings.Trim(nameAndValue[0], " ")
		value := strings.Trim(nameAndValue[1], " ")
		parameterToValue[parameter] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return parameterToValue, nil
}

func buildConfiguration(config *Configuration, paramsToValues *map[string]string) error {
	for k, v := range *paramsToValues {
		switch k {
		case "RemoteURL":
			config.RemoteURL = v
		case "DotFilesDir":
			config.DotFilesDir = v
		case "UpdateIntervalSec":
			if v_num, err := strconv.Atoi(v); err == nil {
				config.UpdateIntervalSec = v_num
			} else {
				return err
			}
		default:
			return errors.New("malformed parameter naming in configuration for: " + k)
		}
	}
	return nil
}