// https://www.freedesktop.org/wiki/Software/systemd/dbus/

package main

import (
	"log"
	"os/exec"
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
		log.Fatalf("[ERROR] %v", err)
	}

	chEvents, chErr := s.conn.SubscribeUnits(time.Duration(s.cfg.SubscribeTimeout) * time.Second)

	go func() {
		for {
			select {
			case events := <-chEvents:
				for _, unit := range events {
					s.processEvent(unit)
				}
			case err := <-chErr:
				log.Printf("[ERROR] %v", err)
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

	log.Printf("[INFO] match unit %v, ActiveState %v, SubState %v", u.Name, u.ActiveState, u.SubState)

	switch u.ActiveState {
	case "active":
		s.execute(unit.OnActive)
	case "inactive":
		s.execute(unit.OnActive)
	case "failed":
		s.execute(unit.OnActive)
	}
}

func (s *subscriber) execute(cmds []string) {
	for _, c := range cmds {
		log.Printf("[INFO] execute %v", c)
		cc, err := shlex.Split(c)
		if err != nil {
			log.Printf("[ERROR] %v", err)
		}

		command := exec.Command(cc[0])
		if len(cc) > 1 {
			command = exec.Command(cc[0], cc[1:]...)
		}
		err = command.Start()
		if err != nil {
			log.Printf("[ERROR] %v", err)
			continue
		}
	}
}
