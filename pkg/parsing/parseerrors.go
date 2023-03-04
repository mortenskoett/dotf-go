package parsing

import "fmt"

/* Parse errors */

type ParseHelpFlagError struct {
	message string
}

type ParseNoArgumentError struct {
	message string
}

type ParseError struct {
	message string
}

type ParseInvalidArgumentError struct {
	message string
}

type ParseConfigurationError struct {
	message string
}

func (e *ParseConfigurationError) Error() string {
	return fmt.Sprintf(e.message)
}

func (e *ParseHelpFlagError) Error() string {
	return fmt.Sprintf(e.message)
}

func (e *ParseNoArgumentError) Error() string {
	return fmt.Sprintf(e.message)
}

func (e *ParseError) Error() string {
	return fmt.Sprintf(e.message)
}

func (e *ParseInvalidArgumentError) Error() string {
	return fmt.Sprintf(e.message)
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
