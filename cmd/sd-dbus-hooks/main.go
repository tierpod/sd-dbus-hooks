// https://www.freedesktop.org/wiki/Software/systemd/dbus/

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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

	conn, err := dbus.New()
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	//http.Handle("/unit/start/", unitStartHandler)
	//http.Handle("/unit/stop/", unitStopHandler)
	http.Handle("/unit/status/", unitStatusHandler{conn})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Printf("[INFO] starting web server on: %v\n", cfg.HTTP.Bind)
	err = http.ListenAndServe(cfg.HTTP.Bind, nil)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	// units, err := conn.ListUnits()
	// if err != nil {
	// 	log.Fatalf("[ERROR] %v", err)
	// }
	// for _, u := range units {
	// 	fmt.Printf("%+v\n", u)
	// }

	// err = conn.Subscribe()
	// if err != nil {
	// 	log.Fatalf("[ERROR] %v", err)
	// }

	// chUnits, chErr := conn.SubscribeUnits(time.Second * 5)

	// for {
	// 	select {
	// 	case uu := <-chUnits:
	// 		// &{Name:rsyslog.service Description:System Logging Service LoadState:loaded ActiveState:inactive SubState:dead Followed: Path:/org/freedesktop/systemd1/unit/rsyslog_2eservice JobId:0 JobType: JobPath:/}
	// 		// &{Name:rsyslog.service Description:System Logging Service LoadState:loaded ActiveState:active SubState:running Followed: Path:/org/freedesktop/systemd1/unit/rsyslog_2eservice JobId:0 JobType: JobPath:/}
	// 		for _, u := range uu {
	// 			fmt.Printf("%+v\n", u)
	// 		}
	// 	case err := <-chErr:
	// 		log.Printf("[ERROR] %v", err)
	// 	}
	// }
}
