package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-systemd/dbus"
	"github.com/coreos/go-systemd/sdjournal"
)

const journalNumEntries = 20

type unitJournalHandler struct {
	conn *dbus.Conn
	cfg  *Config
}

func (h unitJournalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/unit/journal/"):]

	jcfg := sdjournal.JournalReaderConfig{
		NumFromTail: 20,
		Matches: []sdjournal.Match{
			{
				Field: sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT,
				Value: name,
			},
		},
	}

	jr, err := sdjournal.NewJournalReader(jcfg)
	if err != nil {
		log.Printf("[ERROR] journal: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("[INFO] journal: show last %v entries for %v", journalNumEntries, name)
	scanner := bufio.NewScanner(jr)
	for scanner.Scan() {
		fmt.Fprintf(w, "%v\n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Printf("[ERROR] journal: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}