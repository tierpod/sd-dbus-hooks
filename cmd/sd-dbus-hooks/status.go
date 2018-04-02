package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-systemd/dbus"
)

type unitStatusHandler struct {
	conn *dbus.Conn
	cfg  *Config
}

func (h unitStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/unit/status/"):]

	// check if unit in config
	unitCfg, err := h.cfg.getUnit(name)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	log.Printf("[INFO] get unit status %v", name)

	units, err := h.conn.ListUnits()
	if err != nil {
		log.Printf("[ERROR] %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, unit := range units {
		if unit.Name == name {
			fmt.Fprintf(w, "Name         %v\n", unit.Name)
			fmt.Fprintf(w, "Description  %v\n", unit.Description)
			fmt.Fprintf(w, "LoadState    %v\n", unit.LoadState)
			fmt.Fprintf(w, "ActiveState  %v\n", unit.ActiveState)
			fmt.Fprintf(w, "SubState     %v\n", unit.ActiveState)
			fmt.Fprintf(w, "BlockedBy    %v\n", unitCfg.BlockedBy)
			return
		}
	}

	log.Printf("[ERROR] unit %v not found", name)
	w.WriteHeader(http.StatusBadRequest)
	return
}
