package ui

import (
    "log"
    "os/exec"
)

func ShowErrorDialogLinux(message string) {
    log.Println("[UI] Error:", message)
    cmd := exec.Command("zenity", "--error", "--text", message)
    err := cmd.Run()
    if err != nil {
        log.Println("[UI] Failed to display error dialog:", err)
    }
}

func ShowInfoDialogLinux(message string) {
    log.Println("[UI] Info:", message)
    cmd := exec.Command("zenity", "--info", "--text", message)
    err := cmd.Run()
    if err != nil {
        log.Println("[UI] Failed to display info dialog:", err)
    }
}
