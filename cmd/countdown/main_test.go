package main

import (
	"testing"
	"time"

	ui "github.com/gizak/termui/v3"
)

func TestSec2Human(t *testing.T) {
	tests := []struct {
		in   int
		want string
	}{
		{-5, "None"},
		{0, "None"},
		{1, "1 seconds"},
		{59, "59 seconds"},
		{60, "01:00"},
		{61, "01:01"},
		{3599, "59:59"},
		{3600, "01:00:00"},
		{3661, "01:01:01"},
		{7200, "02:00:00"},
	}

	for _, tt := range tests {
		if got := sec2Human(tt.in); got != tt.want {
			t.Errorf("sec2Human(%d) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestGetTimes(t *testing.T) {
	now := int(time.Now().Unix())

	percent, human := getTimes(now, 100)
	if percent < 99 || percent > 100 {
		t.Errorf("getTimes() percent = %d at the start, want ~100", percent)
	}
	if human == "None" {
		t.Errorf("getTimes() human = %q at the start", human)
	}

	percent, human = getTimes(now-200, 100)
	if percent >= 1 || human != "None" {
		t.Errorf("getTimes() after expiry = (%d, %q), want (<1, \"None\")", percent, human)
	}
}

func TestBarColor(t *testing.T) {
	tests := []struct {
		percent int
		want    ui.Color
	}{
		{100, ui.ColorGreen},
		{66, ui.ColorGreen},
		{65, ui.ColorYellow},
		{34, ui.ColorYellow},
		{33, ui.ColorRed},
		{0, ui.ColorRed},
	}

	for _, tt := range tests {
		if got := barColor(tt.percent); got != tt.want {
			t.Errorf("barColor(%d) = %v, want %v", tt.percent, got, tt.want)
		}
	}
}
