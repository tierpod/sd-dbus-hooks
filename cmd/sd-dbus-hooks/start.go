package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-systemd/dbus"
)

type unitStartHandler struct {
	conn *dbus.Conn
	cfg  *Config
}

func (h unitStartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/unit/start/"):]

	// check if unit in config
	unit, err := h.cfg.getUnit(name)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := make(chan string)
	log.Printf("[INFO] starting unit %v", name)

	err = start(h.conn, h.cfg, unit, result)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	status := <-result
	switch status {
	case "done":
		log.Printf("[INFO] unit %v started successfull", name)
		return
	case "timeout", "failed":
		log.Printf("[ERROR] unit %v not started: %v", name, status)
		http.Error(w, status, http.StatusInternalServerError)
		return
	}
	return
}

func start(conn *dbus.Conn, cfg *Config, u Unit, ch chan<- string) error {
	// check if unit active or activating
	// units, err := conn.ListUnitsByPatterns([]string{"active"}, []string{u.Name})
	units, err := listUnitsByPatterns(conn, []string{sdStateActive, sdStateActivating}, []string{u.Name})
	if err != nil {
		return err
	}

	if len(units) != 0 {
		return fmt.Errorf("unit %v already active", u.Name)
	}

	if len(u.BlockedBy) > 0 {
		// check if unit blocked by other active or activating unit
		// blockUnits, err := conn.ListUnitsByPatterns([]string{"active"}, u.BlockedBy)
		blockUnits, err := listUnitsByPatterns(conn, []string{sdStateActive, sdStateActivating}, u.BlockedBy)
		if err != nil {
			return err
		}

		if len(blockUnits) != 0 {
			return fmt.Errorf("unit %v blocked by active units %+v", u.Name, blockUnits)
		}
	}

	// start unit
	_, err = conn.StartUnit(u.Name, "fail", ch)
	if err != nil {
		return err
	}

	return nil
}
