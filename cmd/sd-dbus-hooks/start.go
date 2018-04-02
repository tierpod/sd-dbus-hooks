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

	result := make(chan string)
	log.Printf("[INFO] starting unit %v", name)

	err := start(h.conn, h.cfg, name, result)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status := <-result
	switch status {
	case "done":
		log.Printf("[INFO] unit %v started successfull", name)
		return
	case "timeout":
		log.Printf("[ERROR] unit %v not started: timeout error", name)
		w.WriteHeader(http.StatusBadRequest)
		return
	case "failed":
		log.Printf("[ERROR] unot %v not started: failed", name)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	return
}

func start(conn *dbus.Conn, cfg *Config, name string, ch chan<- string) error {
	// check if unit in config
	unit, err := cfg.getUnit(name)
	if err != nil {
		return err
	}

	// check if unit active
	units, err := conn.ListUnitsByPatterns([]string{"active"}, []string{name})
	if err != nil {
		return err
	}

	if len(units) != 0 {
		return fmt.Errorf("unit %v already active", name)
	}

	// check if unit blocked by other active unit
	blockUnits, err := conn.ListUnitsByPatterns([]string{"active"}, unit.BlockedBy)
	if err != nil {
		return err
	}

	if len(blockUnits) != 0 {
		return fmt.Errorf("unit %v blocked by active units %+v", name, blockUnits)
	}

	// start unit
	_, err = conn.StartUnit(name, "fail", ch)
	if err != nil {
		return err
	}

	return nil
}
