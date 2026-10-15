package main

import (
	"fmt"
	"time"
)

// formatWeek renders t as wk<yy><ww>.<d>: the two-digit ISO year, the ISO
// week number and the day of the week (Sunday is 0).
func formatWeek(t time.Time) string {
	year, week := t.ISOWeek()
	year = year % 100
	dow := int(t.Weekday())
	return fmt.Sprintf("wk%02d%02d.%d", year, week, dow)
}

func main() {
	fmt.Println(formatWeek(time.Now()))
}

// vim: noexpandtab filetype=go
