package main

import (
	"testing"
	"time"
)

func TestFormatWeek(t *testing.T) {
	tests := []struct {
		date time.Time
		want string
	}{
		{time.Date(2026, 10, 15, 12, 0, 0, 0, time.UTC), "wk2642.4"}, // Thursday
		{time.Date(2026, 10, 18, 12, 0, 0, 0, time.UTC), "wk2642.0"}, // Sunday ends ISO week 42
		{time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC), "wk2653.5"},   // belongs to ISO year 2026
		{time.Date(2030, 1, 7, 12, 0, 0, 0, time.UTC), "wk3002.1"},
	}

	for _, tt := range tests {
		if got := formatWeek(tt.date); got != tt.want {
			t.Errorf("formatWeek(%s) = %q, want %q", tt.date.Format("2006-01-02"), got, tt.want)
		}
	}
}
