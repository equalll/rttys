package main

import (
	"os"
	"runtime"
	"time"

	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"github.com/zhaojh329/rttys/version"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		return
	}
	log.AddHook(lfshook.NewHook("/var/log/rttys.log", &log.TextFormatter{}))
}

func main() {
	cfg := parseConfig()

	if cfg.httpUsername == "" {
		log.Fatal("You must configure the http username by commandline or config file")
	}

	log.Info("Go Version: ", runtime.Version())
	log.Info("Go OS/Arch: ", runtime.GOOS, "/", runtime.GOARCH)

	log.Info("Rttys Version: ", version.Version())

	gitCommit := version.GitCommit()
	buildTime := version.BuildTime()

	if gitCommit != "" {
		log.Info("Git Commit: ", version.GitCommit())
	}

	if buildTime != "" {
		log.Info("Build Time: ", version.BuildTime())
	}

	br := newBroker(cfg.token)
	go br.run()

	go listenDevice(br, cfg)
	go httpStart(br, cfg)

	for {
		time.Sleep(time.Second)
	}
}
