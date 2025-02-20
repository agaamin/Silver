package ui

import (
    "log"
    "github.com/getlantern/systray"
)

var syncMenuItem *systray.MenuItem

func InitTray() {
    systray.Run(onTrayReady, onTrayExit)
}

func onTrayReady() {
    systray.SetTitle("SPV Resolver")
    syncMenuItem = systray.AddMenuItem("Syncing...", "Current sync status")
    go handleTrayEvents()
}

func onTrayExit() {
    log.Println("[Tray] Exiting system tray...")
}

func handleTrayEvents() {
    for {
        select {
        case <-syncMenuItem.ClickedCh:
            log.Println("[Tray] Manual sync triggered from tray.")
            // Trigger sync logic here
        }
    }
}

func UpdateTrayStatus(status string) {
    log.Println("[Tray] Updating sync status:", status)
    syncMenuItem.SetTitle(status)
}
