package main

import (
	"encoding/json"
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
	_, err := h.cfg.getUnit(name)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	log.Printf("[INFO] get unit status %v", name)

	// units, err := h.conn.ListUnits()
	units, err := h.conn.ListUnitsByPatterns([]string{"active", "inactive", "failed"}, []string{name})
	if err != nil {
		log.Printf("[ERROR] %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(units)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
