/*
Runs a daemon that listens for changes on a designated remote.
*/
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mortenskoett/dotf-go/pkg/worker"
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
	// conf, err := tomlparser.ReadConfigurationFile(projectpath.Root + "/config.toml")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // For testing.
	// conf.DotFilesDir = "/home/mskk/Repos/temp/git/example1"

	// /* Testing github module*/
	// success, err := terminalio.SyncLocalAndRemote(conf.DotFilesDir)
	// if err != nil {
	// 	// log.Fatal(err)
	// 	switch e := err.(type) {
	// 	case *terminalio.ShellExecError:
	// 		fmt.Println("ShellExecError!!!" + e.Error())
	// 	case *terminalio.MergeFailError:
	// 		fmt.Println("MergeFailerror!!!" + e.Error())
	// 	default:
	// 		fmt.Println(e)
	// 	}
	// }

	// if success {
	// 	fmt.Println("Mission accomplished:", success)
	// } else {
	// 	fmt.Println("Mission FAILED:", success)
	// }

	/* worker */
	w := worker.NewWorkerParam(time.Second*2, func() {
		fmt.Println("hello")
		// time.Sleep(2 * time.Second)
	})

	w.Start()

	time.Sleep(5 * time.Hour)

	// w.Stop()
}
