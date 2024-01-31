package bblocks

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func FormatSize(size int64) string {
	if size > 1e9 {
		return fmt.Sprintf("%dGB", size/(1024*1024*1024))
	} else if size > 1e6 {
		return fmt.Sprintf("%dMB", size/(1024*1024))
	} else {
		return fmt.Sprintf("%dKB", size/1024)
	}
}

func NewLimiter(downloadSpeed int) *RateLimiter {
	return &RateLimiter{
		limit:  float64(downloadSpeed),
		burst:  1024 * 1024, // Allow bursts up to 1 MB
		tokens: 0,
	}
}

// Limit returns the allowed download speed in bytes per second.
func (r *RateLimiter) Limit() float64 {
	return r.limit
}

// Reserve calculates the time required to download a given number of bytes.
func (r *RateLimiter) Reserve(bytes int) time.Duration {
	requiredTokens := float64(bytes)
	timeRequired := time.Duration(requiredTokens / r.limit * float64(time.Second))
	return timeRequired
}

func ParseRateLimit(rateLimit string) (int, error) {
	numericPart := rateLimit[:len(rateLimit)-1]
	suffix := rateLimit[len(rateLimit)-1]

	value, err := strconv.ParseFloat(numericPart, 64)
	if err != nil {
		return 0, err
	}

	switch strings.ToUpper(string(suffix)) {
	case "K":
		value *= 1024
	case "M":
		value *= 1024 * 1024
	}

	return int(value), nil
}
