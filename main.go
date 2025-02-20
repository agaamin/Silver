package main

import (
	"context"
	"errors"
	"silverbullet/internal/config"
	"silverbullet/internal/config/auto"
	"silverbullet/internal/resolvers"
	"silverbullet/internal/resolvers/spvproc"
	"flag"

	"silverbullet/internal/ui"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/emersion/go-autostart"
	"github.com/pkg/browser"
	"github.com/randomlogin/sane"
	"github.com/randomlogin/sane/resolver"
)

const Version = "0.0.4-beta2"

type App struct {
	proc             *spvproc.SPVProc // Replacing HNSD with SPV
	server           *http.Server
	config           *config.App
	usrConfig        *config.User
	proxyURL         string
	autostart        *autostart.App
	autostartEnabled bool
	cancel           func()
}

func setupApp() *App {
	c, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	c.Version = Version
	c.DNSProcPath = path.Join(c.Path, "libs", "spvproc.exe") // Updated for SPV
	app, err := NewApp(c)
	if err != nil {
		log.Fatal(err)
	}

	return app
}

func (app *App) setRecursiveAddress() {
	app.usrConfig.RecursiveAddr = config.DefaultDOHUrl
	serv, err := app.newProxyServer()
	if err != nil {
		log.Fatal(err)
	}
	app.server = serv
	var spvProc *spvproc.SPVProc
	if spvProc, err = spvproc.NewSPVProc(app.config.DNSProcPath); err != nil {
		log.Fatal(err)
	}
	app.proc = spvProc
}

func main() {
	showVersion := flag.Bool("version", false, "Print the version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s
", Version)
		os.Exit(0)
	}

	app := setupApp()
	serverErrCh := make(chan error)

	uiconf := func() {
		ui.Data.SetOptionsEnabled(true)
		ui.Data.SetStarted(true)
		ctx, cancel := context.WithCancel(context.Background())
		app.cancel = cancel
		go func() { serverErrCh <- app.listen() }()
	}

	ui.OnStart = uiconf
	ui.OnExit = func() {
		app.stop()
	}

	ui.Loop()
}
