// https://www.freedesktop.org/wiki/Software/systemd/dbus/

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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

	tokens := newTokenStore(cfg.HTTP.XToken)

	conn, err := dbus.New()
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	http.Handle("/unit/start/", tokens.middleware(unitStartHandler{conn, cfg}))
	http.Handle("/unit/stop/", tokens.middleware(unitStopHandler{conn, cfg}))
	http.Handle("/unit/status/", tokens.middleware(unitStatusHandler{conn, cfg}))
	http.Handle("/unit/journal/", tokens.middleware(unitJournalHandler{conn, cfg}))

	if _, err := os.Stat("webui"); os.IsNotExist(err) {
		log.Printf("[WARN] webui directory is not exist, skip starting webui")
	} else {
		log.Printf("[INFO] starting webui")
		http.Handle("/webui/", http.StripPrefix("/webui/", http.FileServer(http.Dir("webui"))))
	}

	log.Printf("[INFO] subscribe to systemd events with interval %v\n", cfg.SubscribeInterval)
	s := newSubscriber(conn, cfg)
	s.subscribe()

	daemon.SdNotify(false, daemon.SdNotifyReady)

	log.Printf("[INFO] starting web server on: %v (%v)\n", cfg.HTTP.Bind, version)
	err = http.ListenAndServe(cfg.HTTP.Bind, nil)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
}
