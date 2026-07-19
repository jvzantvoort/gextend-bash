package main

import "testing"

func TestConsoleWidth(t *testing.T) {
	t.Run("falls back to COLUMNS env var", func(t *testing.T) {
		t.Setenv("COLUMNS", "123")
		if cols, _ := getWinsize(); cols > 0 {
			t.Skip("skipping: a real terminal is attached, COLUMNS fallback is not exercised")
		}
		if got := consoleWidth(); got != 123 {
			t.Errorf("consoleWidth() = %d, want %d", got, 123)
		}
	})

	t.Run("falls back to 80 when nothing is set", func(t *testing.T) {
		if cols, _ := getWinsize(); cols > 0 {
			t.Skip("skipping: a real terminal is attached, default fallback is not exercised")
		}
		t.Setenv("COLUMNS", "")
		if got := consoleWidth(); got != 80 {
			t.Errorf("consoleWidth() = %d, want %d", got, 80)
		}
	})
}

func TestConsoleHeight(t *testing.T) {
	t.Run("falls back to LINES env var", func(t *testing.T) {
		t.Setenv("LINES", "45")
		if _, rows := getWinsize(); rows > 0 {
			t.Skip("skipping: a real terminal is attached, LINES fallback is not exercised")
		}
		if got := consoleHeight(); got != 45 {
			t.Errorf("consoleHeight() = %d, want %d", got, 45)
		}
	})

	t.Run("falls back to 24 when nothing is set", func(t *testing.T) {
		if _, rows := getWinsize(); rows > 0 {
			t.Skip("skipping: a real terminal is attached, default fallback is not exercised")
		}
		t.Setenv("LINES", "")
		if got := consoleHeight(); got != 24 {
			t.Errorf("consoleHeight() = %d, want %d", got, 24)
		}
	})
}
