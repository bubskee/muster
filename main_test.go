package main

import (
	"testing"

	"github.com/bubskee/muster/engine"
)

func TestSelectControls(t *testing.T) {
	controls := []engine.Control{
		engine.EventControl{ID: "v1"},
		engine.EventControl{ID: "p1"},
		engine.EventControl{ID: "r1"},
	}

	got, err := selectControls(controls, "v1,r1")
	if err != nil {
		t.Fatal(err)
	}

	ids := controlIDs(got)
	if len(ids) != 2 || ids[0] != "v1" || ids[1] != "r1" {
		t.Fatalf("unexpected selected controls: %#v", ids)
	}
}

func TestSelectControlsRejectsUnknownID(t *testing.T) {
	controls := []engine.Control{engine.EventControl{ID: "v1"}}
	if _, err := selectControls(controls, "missing"); err == nil {
		t.Fatal("expected unknown control id error")
	}
}

func TestSelectControlsEmptyMeansAll(t *testing.T) {
	controls := []engine.Control{
		engine.EventControl{ID: "v1"},
		engine.EventControl{ID: "r1"},
	}

	got, err := selectControls(controls, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(controls) {
		t.Fatalf("got %d controls, want %d", len(got), len(controls))
	}
}
