package ui

import (
    "log"
    "github.com/getlantern/systray"
)

var syncStatusItem *systray.MenuItem

func InitUI() {
    systray.Run(onReady, onExit)
}

func onReady() {
    systray.SetTitle("SPV Sync")
    syncStatusItem = systray.AddMenuItem("Syncing...", "Current sync status")
}

func onExit() {
    log.Println("[UI] System tray exiting...")
}

func ShowSyncStatus(status string) {
    log.Println("[UI] Sync Status:", status)
    systray.SetTitle("SPV Sync: " + status)
    syncStatusItem.SetTitle(status)
}

func ShowErrorDialog(message string) {
    log.Println("[UI] Error:", message)
    // Placeholder for an actual UI pop-up in a GUI environment
}
