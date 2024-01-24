package bblocks

import (
	"fmt"
	"time"
)

func DisplayDate(start bool) {
	dt := time.Now()

	if start {
		fmt.Printf("start at %s\n", dt.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Printf("finished at %s\n", dt.Format("2006-01-02 15:04:05"))

	}
}
