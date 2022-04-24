// This package contains functions that wrap the default logger interface and enhance it with
// coloring for the different logging levels.
package logger

import (
	"log"
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
	logWithColor(Default, "", str...)
}

func LogSuccess(str ...interface{}) {
	logWithColor(Green, ok, str...)
}

func LogWarn(str ...interface{}) {
	logWithColor(Yellow, warn, str...)
}

func LogError(str ...interface{}) {
	logWithColor(Red, error, str...)
}

// Logs and exits program
func LogFatal(str ...interface{}) {
	log.SetPrefix(Color(fatal, Red))
	log.Fatalln(str...)
}

func logWithColor(color TerminalColor, prefix string, str ...interface{}) {
	log.SetPrefix(Color(prefix, color))
	log.Println(str...)
	// log.SetPrefix(programprefix)
}
