/*
Runs a daemon that listens for changes on a designated remote.
*/
package main

import (
	"log"
	"os"
)

func init() {
	log.SetPrefix("dotf-cli: ")
}

type ColorANSI string

const (
	ColorBlack  ColorANSI = "\u001b[30m"
	ColorRed              = "\u001b[31m"
	ColorGreen            = "\u001b[32m"
	ColorYellow           = "\u001b[33m"
	ColorBlue             = "\u001b[34m"
	ColorReset            = "\u001b[0m"
)

func colorize(message string, color ColorANSI) string {
	return string(color) + message + string(ColorReset)
}

func main() {
	numargs := len(os.Args)
	args := os.Args[1:]

	if numargs > 1 {
		handleInput(args)
	}
}

func handleInput(args []string) {
	// Commandline arg parser
	// Must parse:
	// command 	arg1
	// add			path_to_file

}
