package main

import (
	"fmt"
	"log"

	"mskk/dotf-go/pkg/resources"
	"mskk/dotf-go/pkg/systray"
)

func init() {
	log.SetPrefix("trayicon: ")
}

var icon = make([]byte, 0)

/* Starts a running process by putting an icon in the systray. */
func main() {
	bytes, err := resources.Get("icons/d_pink_lower_case.png")
	if err != nil {
		log.Fatal(err)
	}
	icon = bytes

	systray.Run(onReady, onExit)
}

func onExit() {
	fmt.Println("dotf tray manager shutdown.")
}

func onReady() {
	fmt.Print("dotf tray manager starting up.")
	systray.SetTemplateIcon(icon, icon)
	systray.SetTitle("Dotf Tray Manager")
	systray.SetTooltip("Dotf Manager")

	mQuitOrig := systray.AddMenuItem("Quit", "Quit dotf tray manager")
	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
	}()
}
