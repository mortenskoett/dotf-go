package logger

import (
	"fmt"
)

type TerminalColor int

const (
	Red TerminalColor = iota
	Green
	Yellow
	Blue
	Default
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorReset  = "\033[0m"
)

// Color given string and a TerminalColor, will insert color codes that are
// interpreted in the terminal as color. The color is reset afterwards.
func Color(text string, color TerminalColor) string {
	return colorCode(color) + text + string(colorReset)
}

func colorCode(code TerminalColor) string {
	colors := []string{
		colorRed,
		colorGreen,
		colorYellow,
		colorBlue,
		colorReset, // used e.g. to reset logger
	}

	if int(code) > len(colors) {
		fmt.Println("colorCode error: TerminalColor out of bounds.")
		return colorReset
	}

	return colors[code]
}
