/*
Starts a running process as an icon in the systray.
*/
package main

import (
	"os"
	"time"

	"github.com/mortenskoett/dotf-go/pkg/argparse"
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
	programVersion string = "" // Inserted by build process
	programName    string = "dotf-tray"
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
	logging.Info("Starting", programName, "service.")
	latestReadConf = *readConfiguration()
	updateWorker = *concurrency.NewIntervalWorkerParam(
		time.Second*time.Duration(latestReadConf.UpdateIntervalSec), handleUpdateNowEvent)
	systray.Run(onReady, onExit)
	logging.Info(programName, "service stopped")
}

func readConfiguration() *config.DotfConfiguration {
	args := os.Args[1:] // Ignore executable name
	vflags := argparse.ValueFlags([]string{"config"})

	flags, err := argparse.ParseFlags(args, vflags)
	if err != nil {
		logging.Fatal(err)
	}

	conf, err := argparse.ParseDotfConfig(flags)
	if err != nil {
		logging.Fatal(err)
	}
	logging.Info("Configuration successfully read")
	return conf
}

func onExit() {
	logging.Info(programName, "shutting down")
}

// Main event loop.
func onReady() {
	systray.SetTitle(programName)
	systray.SetTemplateIcon(getDefaultIcon())
	mLastUpdated.Disable()
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
		logging.Info("Toggle auto-update OFF.")
		mToggleUpdate.Uncheck()
		updateWorker.Stop()

	} else {
		logging.Info("Toggle auto-update ON.")
		mToggleUpdate.Check()
		updateWorker.Start()
	}
}

func handleUpdateNowEvent() {
	logging.Info("Updating now")
	systray.SetTemplateIcon(getLoadingIcon())

	err := terminalio.SyncLocalRemote(latestReadConf.DotfilesDir)
	if err != nil {
		showError(err.Error())
		return
	}

	lastUpdated = time.Now().Format(time.Stamp)
	systray.SetTemplateIcon(getDefaultIcon())
	logging.Info("Updating done")
}

func showError(err string) {
	logging.Info(err)
	systray.SetTemplateIcon(getErrorIcon())
	mError.SetTitle(err)
	mError.Show()
}

func getDefaultIcon() []byte {
	bytes, err := resource.GetIcon(resource.PinkLowerCase)
	if err != nil {
		logging.Fatal(err)
	}
	return bytes
}

func getLoadingIcon() []byte {
	bytes, err := resource.GetIcon(resource.PinkLowerCaseTimeGlass)
	if err != nil {
		logging.Fatal(err)
	}
	return bytes
}

func getErrorIcon() []byte {
	// TODO: Change icon to white exclamation mark
	bytes, err := resource.GetIcon(resource.PinkLowerCaseDragon)
	if err != nil {
		logging.Fatal(err)
	}
	return bytes
}
