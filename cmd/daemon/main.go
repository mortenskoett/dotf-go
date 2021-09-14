/*
Runs a daemon that listens for changes on a designated remote.
*/
package main

import (
	"fmt"
	"log"

	"github.com/mortenskoett/dotf-go/pkg/projectpath"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
	"github.com/mortenskoett/dotf-go/pkg/tomlparser"
)

func init() {
	log.SetPrefix("daemon: ")
}

/* Currently this daemon is used for testing */
func main() {
	/* Testing infinite loop */
	// for {
	// 	// TODO: Dummy implementation
	// 	fmt.Println("Daemon hello")
	// 	path := projectpath.Root
	// 	fmt.Println("Daemon started from", path)
	// 	time.Sleep(time.Second * 3)
	// }

	/* Testing configuration parser */
	conf, err := tomlparser.ReadConfigurationFile(projectpath.Root + "/config.toml")
	if err != nil {
		log.Fatal(err)
	}

	// For testing.
	conf.DotFilesDir = "/home/mskk/Repos/temp/git/example1"

	/* Testing github module*/
	success, err := terminalio.SyncLocalAndRemote(conf.DotFilesDir)
	if err != nil {
		// log.Fatal(err)
		switch e := err.(type) {
		case *terminalio.ShellExecError:
			fmt.Println("ShellExecError!!!" + e.Error())
		case *terminalio.MergeFailError:
			fmt.Println("MergeFailerror!!!" + e.Error())
		default:
			fmt.Println(e)
		}
	}

	if success {
		fmt.Println("Mission accomplished:", success)
	} else {
		fmt.Println("Mission FAILED:", success)
	}
}
