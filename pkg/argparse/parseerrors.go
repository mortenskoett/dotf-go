package argparse

import "fmt"

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
