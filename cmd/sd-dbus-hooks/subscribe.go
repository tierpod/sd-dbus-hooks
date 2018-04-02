package main

import (
	"fmt"
	"log"
	"time"

	"github.com/coreos/go-systemd/dbus"
)

type subscriber struct {
	conn *dbus.Conn
	cfg  *Config
}

func newSubscriber(conn *dbus.Conn, cfg *Config) *subscriber {
	return &subscriber{conn: conn, cfg: cfg}
}

func (s *subscriber) subscribe() {
	err := s.conn.Subscribe()
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	chEvents, chErr := s.conn.SubscribeUnits(time.Duration(s.cfg.SubscribeTimeout) * time.Second)

	go func() {
		for {
			select {
			case events := <-chEvents:
				// &{Name:rsyslog.service Description:System Logging Service LoadState:loaded ActiveState:inactive SubState:dead Followed: Path:/org/freedesktop/systemd1/unit/rsyslog_2eservice JobId:0 JobType: JobPath:/}
				// &{Name:rsyslog.service Description:System Logging Service LoadState:loaded ActiveState:active SubState:running Followed: Path:/org/freedesktop/systemd1/unit/rsyslog_2eservice JobId:0 JobType: JobPath:/}
				for _, unit := range events {
					fmt.Printf("%+v\n", unit)
				}
			case err := <-chErr:
				log.Printf("[ERROR] %v", err)
			}
		}
	}()
}
