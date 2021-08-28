/*
main runs a daemon that listens for changes on a designated remote.
*/
package main

import (
	"fmt"
	"log"
	"mskk/dotf-go/pkg/projectpath"
	"time"
)

func init() {
	log.SetPrefix("daemon: ")
}

func main() {
	for {
		// TODO: Dummy implementation
		fmt.Println("Daemon hello")
		path := projectpath.Root
		fmt.Println("Daemon started from", path)
		time.Sleep(time.Second * 3)
	}

}
