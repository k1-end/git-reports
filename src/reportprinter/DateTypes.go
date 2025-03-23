package reportprinter

import "time"

type yearData struct {
	Year   int
	Months map[time.Month]map[int]struct {
		Date        time.Time
		CommitCount int
	}
}
