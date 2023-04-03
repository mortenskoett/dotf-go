/*
Starts a running process as an icon in the systray.
*/
package main

import (
	"os"
	"time"

	"github.com/getlantern/systray"

	"github.com/mortenskoett/dotf-go/pkg/concurrency"
	"github.com/mortenskoett/dotf-go/pkg/logging"
	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/resource"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

const logo = `    _       _     __         _		     _  _
 __| | ___ | |_  / _|  ___  | |_  _ _  __ _ | || |
/ _' |/ _ \|  _||  _| |___| |  _|| '_|/ _' | \_. |
\__/_|\___/ \__||_|          \__||_|  \__/_| |__/
`

const (
	programName string = "dotf-tray"
)

// State used by the event loop of the tray icon UI.
var (
	programVersion   string                     = "" // Inserted by build process
	shouldAutoUpdate bool                       = false
	lastUpdated      string                     = "N/A"
	configuration    *parsing.DotfConfiguration = nil                              // Configuration currently loaded.
	updateWorker     concurrency.IntervalWorker = *concurrency.NewIntervalWorker() // Worker handles background updates.
)

// Components registered in order seen in the trayicon dropdown.
var (
	mError        = systray.AddMenuItem("No error.", "If an error happens, it pops up here.")
	mUpdateNow    = systray.AddMenuItem("Update Now", "Pulls latest from remote and pushes changes.")
	mToggleUpdate = systray.AddMenuItemCheckbox("Automatic Updates", "Will at intervals push/pull latest changes.", shouldAutoUpdate)
	mQuit         = systray.AddMenuItem("Quit", "Quit dotf tray manager")
	mLastUpdated  = systray.AddMenuItem("Last Updated: "+lastUpdated, "Time the dotfiles were last updated.")
)

// Global dotf-tray flags
var (
	flagConfig = parsing.NewFlag("config", "path to dotf configuration file")
)

func main() {
	logging.WithColor(logging.Blue, logo)
	logging.Info("Starting", programName, "service.", "Version:", programVersion)

	flags, err := parsing.ParseCommandlineFlags(os.Args[1:])
	if err != nil {
		handleParsingError(err)
	}

	configpath := flags.GetOrEmpty(flagConfig)
	configuration, err = parsing.ParseConfig(configpath)
	if err != nil {
		handleParsingError(err)
	}

	logging.Info("Configuration successfully read")

	if configuration.AutoSync {
		handleToggleUpdateEvent()
	}

	systray.Run(onReady, onExit)
	logging.Info(programName, "service stopped")
}

func handleParsingError(err error) {
	if err != nil {
		switch err.(type) {
		case *parsing.ParseNoArgumentError:
			logging.Warn(err)
		case *parsing.ParseConfigurationError:
			logging.Fatal("failed to parse dotf config:", err)
		default:
			logging.Fatal("unknown parser error:", err)
		}
	}
}

func onExit() {
	logging.Info(programName, "shutting down")
}

// Main event loop.
func onReady() {
	systray.SetTitle(programName)
	systray.SetIcon(getDefaultIcon())
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

	// When being checked to ON by user
	if !mToggleUpdate.Checked() {
		logging.Info("Toggle auto-update ON.")

		updateWorker = *concurrency.NewIntervalWorkerParam(
			time.Second*time.Duration(configuration.SyncIntervalSecs), handleUpdateNowEvent)

		mToggleUpdate.Check()
		updateWorker.Start()
		return
	}

	// When being checked to OFF by user
	if mToggleUpdate.Checked() {
		logging.Info("Toggle auto-update OFF.")
		mToggleUpdate.Uncheck()
		updateWorker.Stop()
		return
	}
}

func handleUpdateNowEvent() {
	logging.Info("Updating now")
	systray.SetIcon(getLoadingIcon())

	err := terminalio.SyncLocalRemote(configuration.SyncDir)
	if err != nil {
		showError(err.Error())
		return
	}

	lastUpdated = time.Now().Format(time.Stamp)

	mLastUpdated.Hide()
	mLastUpdated = nil

	mLastUpdated = systray.AddMenuItem("Last Updated: "+lastUpdated, "Time the dotfiles were last updated.")
	mLastUpdated.Disable()

	systray.SetIcon(getDefaultIcon())
	logging.Info("Updating done")
}

func showError(err string) {
	logging.Info(err)
	systray.SetIcon(getErrorIcon())
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
	bytes,
		err := resource.GetIcon(resource.PinkLowerCaseDragon)
	if err != nil {
		logging.Fatal(err)
	}
	return bytes
}
