package cli

import "fmt"

type CmdErrorSuccess struct {
	message string
}

func (e *CmdErrorSuccess) Error() string {
	return fmt.Sprintf("success: %s", e.message)
}
