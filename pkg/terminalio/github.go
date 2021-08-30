package terminalio

import (
	"fmt"
	"log"
	"os/exec"
)

func status() {
	cmd := exec.Command("git status")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cmd.Stdout)
}

func pull() {
}

func fetch() {
}

func push() {
}
