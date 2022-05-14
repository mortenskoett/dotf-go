/*
Starts a running process visualized as an icon in the systray.
*/
package main

import (
	"log"
	"time"

	"github.com/mortenskoett/dotf-go/pkg/concurrency"
	"github.com/mortenskoett/dotf-go/pkg/config"
	"github.com/mortenskoett/dotf-go/pkg/logging"
	"github.com/mortenskoett/dotf-go/pkg/resource"
	"github.com/mortenskoett/dotf-go/pkg/systray"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

const logo = `    _       _     __         _		     _  _
 __| | ___ | |_  / _|  ___  | |_  _ _  __ _ | || |
/ _' |/ _ \|  _||  _| |___| |  _|| '_|/ _' | \_. |
\__/_|\___/ \__||_|          \__||_|  \__/_| |__/
`

const (
	version            = "" // Inserted by build process
	programName string = "dotf-tray"
)

// State used by the event loop of the tray icon UI.
var (
	shouldAutoUpdate bool                       = false
	lastUpdated      string                     = "N/A"
	latestReadConf   config.DotfConfiguration   = config.NewConfiguration()        // Configuration currently loaded.
	updateWorker     concurrency.IntervalWorker = *concurrency.NewIntervalWorker() // Worker handles background updates.
)

// Components registered in order seen in the trayicon dropdown.
var (
	mLastUpdated  = systray.AddMenuItem("Last Updated: "+lastUpdated, "Time the dotfiles were last updated.")
	mError        = systray.AddMenuItem("No error.", "If an error happens, it pops up here.")
	mUpdateNow    = systray.AddMenuItem("Update Now", "Pulls latest from remote and pushes changes.")
	mToggleUpdate = systray.AddMenuItemCheckbox("Automatic Updates", "Will at intervals push/pull latest changes.", shouldAutoUpdate)
	mQuit         = systray.AddMenuItem("Quit", "Quit dotf tray manager")
)

func main() {
	logging.WithColor(logging.Blue, logo)
	logging.Log("A dotfiles updater tray service.\n")
	log.SetPrefix("trayicon: ")
	latestReadConf = readConfiguration()
	updateWorker = *concurrency.NewIntervalWorkerParam(time.Minute*2, handleUpdateNowEvent)
	systray.Run(onReady, onExit)
}

func readConfiguration() config.DotfConfiguration {
	conf, err := config.ReadFromFile(resource.ProjectRoot + "/config.toml")
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Change this path to either $CONFIG or set specifically using UI
	conf.DotfilesDir = "/home/mskk/Repos/temp/git/example1"
	return conf
}

func onExit() {
	logging.Log("Dotf tray manager shutdown")
}

// Main event loop.
func onReady() {
	logging.Log("Dotf tray manager starting up")
	systray.SetTitle("Dotf Tray Manager")
	systray.SetTemplateIcon(getDefaultIcon())
	mLastUpdated.Disable()
	// mError.Disable()
	mError.Hide()

	// Handle events.
	for {
		select {
		case <-mQuit.ClickedCh:
			systray.Quit()
		case <-mToggleUpdate.ClickedCh:
			handleToggleUpdateEvent()
		case <-mUpdateNow.ClickedCh:
			handleUpdateNowEvent()
		case <-mError.ClickedCh:
			mError.Hide()
		}
	}
}

func handleToggleUpdateEvent() {
	shouldAutoUpdate = !shouldAutoUpdate

	if mToggleUpdate.Checked() {
		log.Println("Toggle auto-update OFF.")
		mToggleUpdate.Uncheck()
		updateWorker.Stop()

	} else {
		log.Println("Toggle auto-update ON.")
		mToggleUpdate.Check()
		updateWorker.Start()
	}
}

func handleUpdateNowEvent() {
	log.Println("Updating now")
	systray.SetTemplateIcon(getLoadingIcon())

	err := terminalio.SyncLocalRemote(latestReadConf.DotfilesDir)
	if err != nil {
		showError(err.Error())
		return
	}

	lastUpdated = time.Now().Format(time.Stamp)
	systray.SetTemplateIcon(getDefaultIcon())
	log.Println("Updating done")
}

func showError(err string) {
	log.Print(err)
	systray.SetTemplateIcon(getErrorIcon())
	mError.SetTitle(err)
	mError.Show()
}

func getDefaultIcon() []byte {
	bytes, err := resource.Get("icons/d_pink_lower_case.png")
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func getLoadingIcon() []byte {
	bytes, err := resource.Get("icons/d_pink_lower_case_timeglass.png")
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func getErrorIcon() []byte {
	// TODO: Change icon to white exclamation mark
	bytes, err := resource.Get("icons/d_pink_lower_case_dragon.png")
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
