package bblocks

import (
	"time"
)

func DisplayDate(start bool) {
	dt := time.Now()

	x := dt.Format("2006-01-02 15:04:05")

	if start {
		outputFunc("start at " + x + "\n")
	} else {
		outputFunc("finished at " + x + "\n")
	}
}
