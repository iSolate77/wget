package bblocks

import (
	"io"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)


func DownloadWithStandardProgressBar(body io.Reader, file *os.File, limiter *RateLimiter, totalSize int64, bar *progressbar.ProgressBar) error {
	var reader io.Reader = body // Use body reader initially

	// Apply rate limiter if provided
	if limiter != nil {
		reader = NewRateLimitedReader(body, limiter) // Wrap body reader with rate limiter
	}

	// Copy data from reader to file with progress bar
	_, err := io.Copy(io.MultiWriter(file, bar), reader)
	if err != nil {
		return err
	}

	return nil
}

func NewRateLimitedReader(r io.Reader, limiter *RateLimiter) io.Reader {
	return &rateLimitedReader{
		reader:  r,
		limiter: limiter,
	}
}

// rateLimitedReader is a rate-limited reader that enforces the rate limit while reading.
type rateLimitedReader struct {
	reader  io.Reader
	limiter *RateLimiter
}

// Read reads data from the underlying reader and enforces the rate limit.
func (rlr *rateLimitedReader) Read(p []byte) (n int, err error) {
	// Reserve tokens based on the size of the buffer
	n, err = rlr.reader.Read(p)
	if err != nil {
		return n, err
	}
	if rlr.limiter != nil {
		// Calculate the time required to read the data and reserve tokens
		timeRequired := rlr.limiter.Reserve(n)
		time.Sleep(timeRequired)
	}
	return n, nil
}
