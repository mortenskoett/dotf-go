package logger

import (
	"log"

	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

const (
	prefix string = "dotf-cli: "
)

func init() {
	log.SetFlags(0)
}

func Log(str interface{}) {
	logWithColor(str, terminalio.Default)
}

func LogSuccess(str interface{}) {
	logWithColor(str, terminalio.Green)
}

func LogWarn(str interface{}) {
	logWithColor(str, terminalio.Yellow)
}

func LogError(str interface{}) {
	logWithColor(str, terminalio.Red)
}

func logWithColor(str interface{}, color terminalio.TerminalColor) {
	log.SetPrefix(terminalio.Color(prefix, color))
	log.Println(str)
	log.SetPrefix(prefix)
}
