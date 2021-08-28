/*
Starts a running process by putting an icon in the systray.
*/
package main

import (
	"fmt"
	"log"
	"time"

	"mskk/dotf-go/pkg/resources"
	"mskk/dotf-go/pkg/systray"
)

func init() {
	log.SetPrefix("trayicon: ")
}

// Setup variables.
var (
	shouldAutoUpdate = false
	lastUpdated      = time.Now().Format(time.Stamp)
)

// Register components in order.
var (
	mLastUpdated  = systray.AddMenuItem("Last Updated: "+lastUpdated, "Time the dotfiles were last updated.")
	mUpdateNow    = systray.AddMenuItem("Update Now", "Pulls latest from remote and pushes changes.")
	mToggleUpdate = systray.AddMenuItemCheckbox("Automatic Updates", "Will at intervals push/pull latest changes.", shouldAutoUpdate)
	mQuit         = systray.AddMenuItem("Quit", "Quit dotf tray manager")
)

func main() {
	systray.Run(onReady, onExit)
}

func onExit() {
	fmt.Println("dotf tray manager shutdown.")
}

// Main event loop.
func onReady() {
	fmt.Print("dotf tray manager starting up.")
	systray.SetTitle("Dotf Tray Manager")
	systray.SetTemplateIcon(getDefaultIcon())
	mLastUpdated.Disable()

	// Handle events.
	for {
		select {
		case <-mQuit.ClickedCh:
			systray.Quit()
		case <-mToggleUpdate.ClickedCh:
			handleToggleUpdateEvent()
		case <-mUpdateNow.ClickedCh:
			handleUpdateNowEvent()
		}
	}
}

func handleToggleUpdateEvent() {
	shouldAutoUpdate = !shouldAutoUpdate
	if mToggleUpdate.Checked() {
		mToggleUpdate.Uncheck()
	} else {
		mToggleUpdate.Check()
	}
}

func handleUpdateNowEvent() {
	//TODO
	// Show a loading icon
	// Push/pull latest dotfiles
	// When operation returns reset icon
	fmt.Println("Updates")
	systray.SetTemplateIcon(getLoadingIcon())
	time.Sleep(time.Second * 10)
	systray.SetTemplateIcon(getDefaultIcon())
}

func getDefaultIcon() []byte {
	bytes, err := resources.Get("icons/d_pink_lower_case.png")
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func getLoadingIcon() []byte {
	bytes, err := resources.Get("icons/d_pink_lower_case_timeglass.png")
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
