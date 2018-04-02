package main

import (
	"log"
	"net/http"

	"github.com/coreos/go-systemd/dbus"
)

type unitStopHandler struct {
	conn *dbus.Conn
}

func (h unitStopHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/unit/stop/"):]
	result := make(chan string)
	log.Printf("[INFO] stopping unit %v", name)

	_, err := h.conn.StopUnit(name, "fail", result)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status := <-result
	switch status {
	case "done":
		log.Printf("[INFO] unit %v stopped successfull", name)
		return
	case "timeout":
		log.Printf("[ERROR] unit %v not stopped: timeout error", name)
		w.WriteHeader(http.StatusBadRequest)
		return
	case "failed":
		log.Printf("[ERROR] unit %v not stopped: failed", name)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	return
}
