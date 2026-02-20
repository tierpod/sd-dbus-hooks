package main

import (
	"path"
	"reflect"

	"github.com/coreos/go-systemd/v22/dbus"
)

func listUnitsByPatterns(conn *dbus.Conn, states []string, patterns []string) ([]dbus.UnitStatus, error) {
	var result []dbus.UnitStatus

	units, err := conn.ListUnits()
	if err != nil {
		return nil, err
	}

	for _, unit := range units {
		if contains(patterns, unit.Name) && contains(states, unit.ActiveState) {
			result = append(result, unit)
			continue
		}
	}

	// do not need to list unit files if we search only "active" and "activating" units
	if !reflect.DeepEqual(states, []string{sdStateActive, sdStateActivating}) {
		// systemd can doesn't show all loaded units in some cases (if there's no reason to keep it in memory)
		// https://github.com/systemd/systemd/issues/5063
		//
		// so, list all units files for matched names and add it to results
		unitFiles, err := conn.ListUnitFiles()
		if err != nil {
			return nil, err
		}

	UNIT_FILES_LOOP:
		for _, unitFile := range unitFiles {
			name := path.Base(unitFile.Path)
			// skip unit file if result already exist
			for _, v := range result {
				if name == v.Name {
					// log.Printf("[DEBUG] unit %v already in results", v.Name)
					continue UNIT_FILES_LOOP
				}
			}
			if contains(patterns, name) {
				// create fake unloaded UnitStatus
				unloaded := dbus.UnitStatus{
					Name:        name,
					Description: "",
					LoadState:   sdStateUnloaded,
					ActiveState: sdStateUnloaded,
					SubState:    unitFile.Type,
				}
				result = append(result, unloaded)
				continue
			}
		}
	}

	return result, nil
}

// return true if string `v` contains in slice `s`
func contains(s []string, v string) bool {
	for _, i := range s {
		if i == v {
			return true
		}
	}

	return false
}
