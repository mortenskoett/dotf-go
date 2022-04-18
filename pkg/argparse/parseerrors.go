package argparse

import "fmt"

type ParseErrorSuccess struct {
	message string
}

func (e *ParseErrorSuccess) Error() string {
	return fmt.Sprintf(e.message)
}
