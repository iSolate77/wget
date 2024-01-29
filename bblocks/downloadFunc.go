package bblocks

import (
	"io"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

func DownloadWithProgressBar(body io.Reader, file *os.File, limiter *RateLimiter, totalSize int64, bar *progressbar.ProgressBar) error {

	var buf []byte
	if totalSize <= 0{
		buf = make([]byte, defaultBufferSize)
	}else{
		buf = make([]byte, totalSize)
	}
	for {
		n, err := body.Read(buf)
		if err != nil {
			if err == io.EOF {
				// If we've reached the end of the response body, break the loop
				break
			}
			return err
		}
		if limiter != nil {
			timeRequired := limiter.Reserve(n)
			time.Sleep(timeRequired)
		}
		_, err = file.Write(buf[:n])
		if err != nil {
			return err
		}
		bar.Add(n)
	}
	return nil
}
