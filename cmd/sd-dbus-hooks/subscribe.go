// https://www.freedesktop.org/wiki/Software/systemd/dbus/

package main

import (
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/coreos/go-systemd/dbus"
	"github.com/google/shlex"
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
		log.Fatalf("[ERROR] subscriber: %v", err)
	}

	chEvents, chErr := s.conn.SubscribeUnits(time.Duration(s.cfg.SubscribeInterval) * time.Second)

	go func() {
		for {
			select {
			case events := <-chEvents:
				for _, unit := range events {
					if unit == nil {
						log.Printf("[WARN] subscriber: got nil event, ignore")
						continue
					}
					s.processEvent(unit)
				}
			case err := <-chErr:
				log.Printf("[ERROR] subscriber: %v", err)
			}
		}
	}()
}

func (s *subscriber) processEvent(u *dbus.UnitStatus) {
	unit, err := s.cfg.getUnit(u.Name)
	if err != nil {
		// log.Printf("[DEBUG] %v", err)
		return
	}

	log.Printf("[INFO] subscriber: match unit %v, ActiveState %v, SubState %v", u.Name, u.ActiveState, u.SubState)

	switch u.ActiveState {
	case "active", "activating":
		go s.execute(unit.OnActive, u)
	case "inactive", "deactivating":
		go s.execute(unit.OnInctive, u)
	case "failed":
		go s.execute(unit.OnFailed, u)
	}
}

func (s *subscriber) execute(cmds []string, u *dbus.UnitStatus) {
	for _, c := range cmds {
		c = strings.Replace(c, "{unit_name}", u.Name, -1)
		c = strings.Replace(c, "{unit_state}", u.ActiveState+"/"+u.SubState, -1)
		log.Printf("[INFO] subscriber: execute %v", c)
		cc, err := shlex.Split(c)
		if err != nil {
			log.Printf("[ERROR] subscriber: %v", err)
		}

		command := exec.Command(cc[0])
		if len(cc) > 1 {
			command = exec.Command(cc[0], cc[1:]...)
		}
		err = command.Run()
		if err != nil {
			log.Printf("[ERROR] subscriber: execute failed: %v", err)
			continue
		}
	}
}
