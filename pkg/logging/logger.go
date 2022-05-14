// This package contains functions that wrap the default logger interface and enhance it with
// coloring for the different logging levels.
package logging

import (
	"log"
)

const (
	warn  string = "warn: "
	fatal string = "fatal: "
	error string = "error: "
	ok    string = "ok: "
	debug string = "DEBUG: "
	info  string = "info: "
)

func init() {
	log.SetFlags(0)
}

func Log(str ...interface{}) {
	logWithColor(Default, "", str...)
}

func Ok(str ...interface{}) {
	logWithColor(Green, ok, str...)
}

func Warn(str ...interface{}) {
	logWithColor(Yellow, warn, str...)
}

func Error(str ...interface{}) {
	logWithColor(Red, error, str...)
}

func Debug(str ...interface{}) {
	logWithColor(Yellow, debug, str...)
}

func Info(str ...interface{}) {
	logWithColor(Blue, info, str...)
}

func WithColor(color TerminalColor, str ...string) {
	log.Println(ColorMultiple(color, str...))
}

// Logs and exits program
func Fatal(str ...interface{}) {
	log.SetPrefix(Color(fatal, Red))
	log.Fatalln(str...)
}

func logWithColor(color TerminalColor, prefix string, str ...interface{}) {
	log.SetPrefix(Color(prefix, color))
	log.Println(str...)
}
