package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	year, week := now.ISOWeek()
	year = int(year) % 2000
	dow := int(now.Weekday())
	fmt.Printf("wk%02d%02d.%d\n", year, week, dow)
}

// vim: noexpandtab filetype=go
