package icon

import (
    "github.com/getlantern/systray"
    "log"
)

func SetTrayIconLinux() {
    log.Println("[Icon] Setting tray icon for Linux.")
    systray.SetIcon(IconLinux)
}

func SetTrayIconWindows() {
    log.Println("[Icon] Setting tray icon for Windows.")
    systray.SetIcon(IconWindows)
}

func SetTrayIconMac() {
    log.Println("[Icon] Setting tray icon for Mac.")
    systray.SetIcon(IconMac)
}
