package main

import (
	"slices"
	"testing"

	tls "github.com/jvzantvoort/gextend-bash/printing"
)

func TestColorNames(t *testing.T) {
	got := colorNames(tls.NewCprint())
	if len(got) == 0 {
		t.Fatal("colorNames() returned no names")
	}
	if !slices.IsSorted(got) {
		t.Errorf("colorNames() = %v, want sorted", got)
	}
	if !slices.Contains(got, "ok") {
		t.Errorf("colorNames() = %v, want it to contain %q", got, "ok")
	}
}
