// This package contains functions that wrap the default logger interface and enhance it with
// coloring for the different logging levels.
package logger

import (
	"log"

	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

const (
	// programprefix string = constant.ProgramName + ": "
	warn  string = "warn: "
	fatal string = "fatal: "
	error string = "error: "
	ok    string = "ok: "
)

func init() {
	log.SetFlags(0)
}

func Log(str ...interface{}) {
	logWithColor(terminalio.Default, "", str...)
}

func LogSuccess(str ...interface{}) {
	logWithColor(terminalio.Green, ok, str...)
}

func LogWarn(str ...interface{}) {
	logWithColor(terminalio.Yellow, warn, str...)
}

func LogError(str ...interface{}) {
	logWithColor(terminalio.Red, error, str...)
}

// Logs and exits program
func LogFatal(str ...interface{}) {
	log.SetPrefix(terminalio.Color(fatal, terminalio.Red))
	log.Fatalln(str...)
}

func logWithColor(color terminalio.TerminalColor, prefix string, str ...interface{}) {
	log.SetPrefix(terminalio.Color(prefix, color))
	log.Println(str...)
	// log.SetPrefix(programprefix)
}
