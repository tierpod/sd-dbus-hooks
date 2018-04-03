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

	cfg, err := LoadConfig(flagConfig)
	if err != nil {
		log.Fatal(err)
	}

	logFlags := 0
	if cfg.HTTP.LogTimestamp {
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
	http.Handle("/units", tokens.middleware(unitsHandler{cfg}))

	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Printf("[INFO] subscribe to systemd events with timeout %v\n", cfg.SubscribeTimeout)
	s := newSubscriber(conn, cfg)
	s.subscribe()

	daemon.SdNotify(false, daemon.SdNotifyReady)

	log.Printf("[INFO] starting web server on: %v\n", cfg.HTTP.Bind)
	err = http.ListenAndServe(cfg.HTTP.Bind, nil)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
}

func contains(s string, ss []string) bool {
	for _, i := range ss {
		if i == s {
			return true
		}
	}

	return false
}
