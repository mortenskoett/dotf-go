/*
Starts a running process by putting an icon in the systray.
*/
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mortenskoett/dotf-go/pkg/projectpath"
	"github.com/mortenskoett/dotf-go/pkg/resources"
	"github.com/mortenskoett/dotf-go/pkg/systray"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
	"github.com/mortenskoett/dotf-go/pkg/tomlparser"
	"github.com/mortenskoett/dotf-go/pkg/worker"
)

func init() {
	log.SetPrefix("trayicon: ")
}

// Setup variables.
var (
	shouldAutoUpdate bool                     = false
	lastUpdated      string                   = "N/A"
	latestReadConf   tomlparser.Configuration = tomlparser.NewConfiguration() // Configuration currently loaded.
	updateWorker     worker.Worker            = *worker.NewWorker()           // Worker handles background updates.
)

// Register components in order seen in the trayicon dropdown.
var (
	mLastUpdated  = systray.AddMenuItem("Last Updated: "+lastUpdated, "Time the dotfiles were last updated.")
	mError        = systray.AddMenuItem("No error.", "If an error happens, it pops up here.")
	mUpdateNow    = systray.AddMenuItem("Update Now", "Pulls latest from remote and pushes changes.")
	mToggleUpdate = systray.AddMenuItemCheckbox("Automatic Updates", "Will at intervals push/pull latest changes.", shouldAutoUpdate)
	mQuit         = systray.AddMenuItem("Quit", "Quit dotf tray manager")
)

func main() {
	latestReadConf = readConfiguration()
	updateWorker = *worker.NewWorkerParam(time.Minute*2, handleUpdateNowEvent)
	systray.Run(onReady, onExit)
}

func readConfiguration() tomlparser.Configuration {
	conf, err := tomlparser.ReadConfigurationFile(projectpath.Root + "/config.toml")
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Change this path to either $CONFIG or set specifically using UI
	conf.DotFilesDir = "/home/mskk/Repos/temp/git/example1"
	return conf
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

	_, err := terminalio.SyncLocalAndRemote(latestReadConf.DotFilesDir)
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

func getErrorIcon() []byte {
	// TODO: Change icon to white exclamation mark
	bytes, err := resources.Get("icons/d_pink_lower_case_dragon.png")
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
