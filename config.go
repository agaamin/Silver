package config

import (
    "bytes"
    "crypto/tls"
    "crypto/x509"
    "encoding/json"
    "encoding/pem"
    "fmt"
    "math/rand"
    "net/http"
    "os"
    "path"
    "strings"
    "time"
)

const (
    AppName          = "SilverBullet"
    AppId            = "com.impervious.silverbullet"
    CertFileName     = "silverbullet.crt"
    CertKeyFileName  = "private.key"
    CertName         = "DNSSEC"
    DefaultSyncInterval = 12 * time.Hour // Default sync interval
)

type App struct {
    Path        string
    CertPath    string
    keyPath     string
    DNSProcPath string
    ProxyAddr   string
    Version     string
    SyncInterval time.Duration
    DLLPaths    map[string]string
}

func NewConfig() (*App, error) {
    var err error
    c := &App{SyncInterval: DefaultSyncInterval}
    if c.Path, err = getOrCreateDir(); err != nil {
        return nil, fmt.Errorf("failed creating config: %v", err)
    }

    // Define DLL paths for local loading (Updated for SPV instead of HNSD)
    c.DLLPaths = map[string]string{
        "libcrypto":  path.Join(c.Path, "libs", "libcrypto.dll"),
        "libssl":     path.Join(c.Path, "libs", "libssl.dll"),
        "libevent":   path.Join(c.Path, "libs", "libevent.dll"),
        "libunbound": path.Join(c.Path, "libs", "libunbound.dll"),
        "spvproc":    path.Join(c.Path, "libs", "spvproc.exe"), // Replacing HNSD with SPV
    }

    return c, nil
}

func (c *App) SetSyncInterval(interval time.Duration) {
    c.SyncInterval = interval
    log.Println("[Config] Sync interval set to:", interval)
}

func (c *App) GetSyncInterval() time.Duration {
    return c.SyncInterval
}

func (c *App) GetDLLPath(name string) string {
    return c.DLLPaths[name]
}
