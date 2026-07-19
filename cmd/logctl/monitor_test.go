package main

import "testing"

func TestMonitorStatePrepend(t *testing.T) {
	m := &monitorState{maxLines: 3}

	m.prepend("a")
	m.prepend("b")
	m.prepend("c")

	want := []string{"c", "b", "a"}
	if len(m.rawLines) != len(want) {
		t.Fatalf("rawLines = %v, want %v", m.rawLines, want)
	}
	for i := range want {
		if m.rawLines[i] != want[i] {
			t.Errorf("rawLines[%d] = %q, want %q", i, m.rawLines[i], want[i])
		}
	}

	t.Run("evicts the oldest entry once at capacity", func(t *testing.T) {
		m.prepend("d")
		want := []string{"d", "c", "b"}
		if len(m.rawLines) != len(want) {
			t.Fatalf("rawLines = %v, want %v", m.rawLines, want)
		}
		for i := range want {
			if m.rawLines[i] != want[i] {
				t.Errorf("rawLines[%d] = %q, want %q", i, m.rawLines[i], want[i])
			}
		}
	})
}

func TestMonitorStateResizeTrimsBuffer(t *testing.T) {
	m := &monitorState{maxLines: 10, rawLines: []string{"a", "b", "c", "d", "e"}}

	if cols, _ := getWinsize(); cols > 0 {
		t.Skip("skipping: a real terminal is attached, resize() would read live dimensions")
	}
	t.Setenv("LINES", "2")
	t.Setenv("COLUMNS", "40")

	m.resize()

	if m.maxLines != 2 {
		t.Fatalf("maxLines = %d, want 2", m.maxLines)
	}
	if len(m.rawLines) != 2 {
		t.Errorf("rawLines = %v, want length 2", m.rawLines)
	}
	if m.width != 40 {
		t.Errorf("width = %d, want 40", m.width)
	}
}
