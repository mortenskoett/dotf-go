package parsing

import "fmt"

/* Config errors */

type ConfigKeyNotFoundError struct {
	message string
}

func (e *ConfigKeyNotFoundError) Error() string {
	return e.message
}

type MalformedConfigurationError struct {
	message string
}

func (e *MalformedConfigurationError) Error() string {
	return e.message
}

/* Parse errors */

type ParseNoArgumentError struct {
	message string
}

type ParseDefaultConfigurationError struct {
	message string
}

type ParseConfigurationError struct {
	message string
}

type ParseInvalidFlagError struct {
	message string
}

func (e *ParseNoArgumentError) Error() string {
	return fmt.Sprintf(e.message)
}

func (e *ParseDefaultConfigurationError) Error() string {
	return fmt.Sprintf("invalid configuration found at default location: %v", e.message)
}

func (e *ParseConfigurationError) Error() string {
	return fmt.Sprintf("failed to parse configuration file: %v", e.message)
}

func (e *ParseInvalidFlagError) Error() string {
	return fmt.Sprintf("failed to parse command line flags: %s", e.message)
}

/* Other errors */

type CombinedError struct {
	Errors []error
}

func (e *CombinedError) Error() string {
	var error string
	for _, err := range e.Errors {
		error += err.Error()
	}
	return error
}
