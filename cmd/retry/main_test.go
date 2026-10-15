package main

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	tests := []struct {
		name       string
		opts       Options
		codes      []int // exit codes returned by successive attempts
		wantCode   int
		wantRuns   int
		wantSleeps []time.Duration
	}{
		{"succeeds first time", Options{Attempts: 3, Delay: time.Second}, []int{0}, 0, 1, nil},
		{"succeeds after failures", Options{Attempts: 3, Delay: time.Second}, []int{1, 1, 0}, 0, 3,
			[]time.Duration{time.Second, time.Second}},
		{"gives up with the last code", Options{Attempts: 3, Delay: time.Second}, []int{1, 2, 3}, 3, 3,
			[]time.Duration{time.Second, time.Second}},
		{"backs off up to the maximum", Options{Attempts: 5, Delay: time.Second, Backoff: true, MaxDelay: 3 * time.Second},
			[]int{1, 1, 1, 1, 1}, 1, 5, []time.Duration{time.Second, 2 * time.Second, 3 * time.Second, 3 * time.Second}},
		{"does not retry a missing command", Options{Attempts: 3, Delay: time.Second}, []int{127}, 127, 1, nil},
		{"single attempt", Options{Attempts: 1, Delay: time.Second}, []int{4}, 4, 1, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runs := 0
			var sleeps []time.Duration
			var log bytes.Buffer
			attempt := func() int { runs++; return tt.codes[runs-1] }
			sleep := func(d time.Duration) { sleeps = append(sleeps, d) }

			code := retry(tt.opts, attempt, sleep, &log)
			if code != tt.wantCode {
				t.Errorf("retry() = %d, want %d", code, tt.wantCode)
			}
			if runs != tt.wantRuns {
				t.Errorf("runs = %d, want %d", runs, tt.wantRuns)
			}
			if !reflect.DeepEqual(sleeps, tt.wantSleeps) {
				t.Errorf("sleeps = %v, want %v", sleeps, tt.wantSleeps)
			}
			if code != 0 && !strings.Contains(log.String(), "giving up") {
				t.Errorf("log = %q, want a giving up message", log.String())
			}
		})
	}
}
