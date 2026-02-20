package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-systemd/dbus"
)

type unitStopHandler struct {
	conn *dbus.Conn
	cfg  *Config
}

func (h unitStopHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/unit/stop/"):]

	// check if unit in config
	_, err := h.cfg.getUnit(name)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := make(chan string)
	log.Printf("[INFO] stopping unit %v", name)

	_, err = h.conn.StopUnit(name, "fail", result)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	status := <-result
	switch status {
	case "done":
		log.Printf("[INFO] unit %v stopped successfull", name)
		fmt.Fprint(w, "OK")
		return
	case "timeout", "failed":
		log.Printf("[ERROR] unit %v not stopped: %v", name, status)
		http.Error(w, status, http.StatusInternalServerError)
		return
	}
}
