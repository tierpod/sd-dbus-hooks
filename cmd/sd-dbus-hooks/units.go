package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type unitsHandler struct {
	cfg *Config
}

func (h unitsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var units []string
	for _, unit := range h.cfg.Units {
		units = append(units, unit.Name)
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
