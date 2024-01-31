package bblocks

import (
	"fmt"
	"math"
)

func RoundBytes(in_bytes int64) string {
	sizeInMB := float64(in_bytes) / (1024 * 1024)

	// Round to gigabytes if size is 1024 MB or larger
	sizeInGB := math.Round(sizeInMB/1024*100) / 100

	if sizeInGB >= 1 {
		return fmt.Sprintf("%.2f GB", sizeInGB)
	} else {
		return fmt.Sprintf("%.2f MB", sizeInMB)
	}
}
