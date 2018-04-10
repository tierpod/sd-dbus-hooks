// https://www.freedesktop.org/wiki/Software/systemd/dbus/

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/go-systemd/daemon"
	"github.com/coreos/go-systemd/dbus"
)

var version string

func main() {
	var (
		flagVersion bool
		flagConfig  string
	)

	flag.BoolVar(&flagVersion, "version", false, "Show version and exit")
	flag.StringVar(&flagConfig, "config", "./config.yaml", "Path to config file")
	flag.Parse()

	if flagVersion {
		fmt.Printf("Version: %v\n", version)
		os.Exit(0)
	}

	cfg, err := loadConfig(flagConfig)
	if err != nil {
		log.Fatal(err)
	}

	logFlags := 0
	if cfg.LogTimestamp {
		logFlags = log.LstdFlags
	}

	log.SetFlags(logFlags)

	sigsCh := make(chan os.Signal, 1)
	shutdownCh := make(chan bool, 1)

	signal.Notify(sigsCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)

	go func() {
		for {
			sig := <-sigsCh
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				log.Printf("[DEBUG] got signal %v, shutdown", sig)
				close(shutdownCh)
			case syscall.SIGUSR2:
				daemon.SdNotify(false, daemon.SdNotifyReloading)
				log.Printf("[INFO] got signal %v, reload config", sig)
				err = updateConfig(cfg, flagConfig)
				if err != nil {
					log.Fatal(err)
				}
				daemon.SdNotify(false, daemon.SdNotifyReady)
			}
		}
	}()

	conn, err := dbus.New()
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	if cfg.HTTP.Bind != "" {
		go startWebServer(conn, cfg)
	}

	log.Printf("[INFO] subscribe to systemd events with interval %v\n", cfg.SubscribeInterval)
	s := newSubscriber(conn, cfg)
	s.subscribe()

	log.Printf("[INFO] service started (version: %v)", version)
	daemon.SdNotify(false, daemon.SdNotifyReady)
	<-shutdownCh
	daemon.SdNotify(false, daemon.SdNotifyStopping)
}

func startWebServer(conn *dbus.Conn, cfg *Config) {
	tokens := tokenStore{cfg: cfg}

	http.Handle("/unit/start/", tokens.middleware(unitStartHandler{conn, cfg}))
	http.Handle("/unit/stop/", tokens.middleware(unitStopHandler{conn, cfg}))
	http.Handle("/unit/status/", tokens.middleware(unitStatusHandler{conn, cfg}))
	http.Handle("/unit/journal/", tokens.middleware(unitJournalHandler{conn, cfg}))

	if _, err := os.Stat("webui"); os.IsNotExist(err) {
		log.Printf("[WARN] webserver: webui directory does not exist, skip starting webui")
	} else {
		log.Printf("[INFO] webserver: starting webui")
		http.Handle("/webui/", http.StripPrefix("/webui/", http.FileServer(http.Dir("webui"))))
	}

	log.Printf("[INFO] webserver: starting web server on: %v", cfg.HTTP.Bind)
	err := http.ListenAndServe(cfg.HTTP.Bind, nil)
	if err != nil {
		log.Fatalf("[ERROR] webserver: %v", err)
	}
}
