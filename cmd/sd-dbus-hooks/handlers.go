package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-systemd/dbus"
)

type unitStatusHandler struct {
	conn *dbus.Conn
}

func (h unitStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/unit/status/"):]
	log.Printf("[INFO] get unit status %v", name)

	units, err := h.conn.ListUnits()
	if err != nil {
		log.Printf("[ERROR] %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, unit := range units {
		if unit.Name == name {
			fmt.Fprintf(w, "Name: %+v\n", unit.Name)
			fmt.Fprintf(w, "Description: %+v\n", unit.Description)
			fmt.Fprintf(w, "LoadState: %+v\n", unit.LoadState)
			fmt.Fprintf(w, "ActiveState: %+v\n", unit.ActiveState)
			fmt.Fprintf(w, "SubState: %+v\n", unit.ActiveState)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	log.Printf("[ERROR] unit %v not found", name)
	w.WriteHeader(http.StatusBadRequest)
	return
}

type unitStartHandler struct {
	conn *dbus.Conn
}

func (h unitStartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/unit/start/"):]
	result := make(chan string)
	log.Printf("[INFO] starting unit %v", name)

	_, err := h.conn.StartUnit(name, "fail", result)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status := <-result
	switch status {
	case "done":
		log.Printf("[INFO] unit %v started successfull", name)
		w.WriteHeader(http.StatusOK)
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
		w.WriteHeader(http.StatusOK)
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
