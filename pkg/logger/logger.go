// This package contains functions that wrap the default logger interface and enhance it with
// coloring for the different logging levels.
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

func Log(str ...interface{}) {
	logWithColor(terminalio.Default, str...)
}

func LogSuccess(str ...interface{}) {
	logWithColor(terminalio.Green, str...)
}

func LogWarn(str ...interface{}) {
	logWithColor(terminalio.Yellow, str...)
}

func LogError(str ...interface{}) {
	logWithColor(terminalio.Red, str...)
}

func LogFatal(str ...interface{}) {
	log.SetPrefix(terminalio.Color(prefix, terminalio.Red))
	log.Fatalln(str...)
}

func logWithColor(color terminalio.TerminalColor, str ...interface{}) {
	log.SetPrefix(terminalio.Color(prefix, color))
	log.Println(str...)
	log.SetPrefix(prefix)
}
