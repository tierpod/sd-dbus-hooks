package main

import "testing"

func TestUtils(t *testing.T) {
	items := []string{"one", "two", "three"}
	if !contains(items, "one") {
		t.Errorf("contains: got false, expected true")
	}

	if contains(items, "four") {
		t.Errorf("contains: got true, expected false")
	}
}
