/*
Runs a daemon that listens for changes on a designated remote.
*/
package main

import (
	"fmt"
	"log"

	"github.com/mortenskoett/dotf-go/pkg/projectpath"
	"github.com/mortenskoett/dotf-go/pkg/tomlparser"
)

func init() {
	log.SetPrefix("daemon: ")
}

func main() {
	// for {
	// 	// TODO: Dummy implementation
	// 	fmt.Println("Daemon hello")
	// 	path := projectpath.Root
	// 	fmt.Println("Daemon started from", path)
	// 	time.Sleep(time.Second * 3)
	// }

	conf, err := tomlparser.ReadConfigurationFile(projectpath.Root + "/config.toml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(conf)
}
