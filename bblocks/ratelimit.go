package bblocks

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var outputFunc func(a ...interface{}) (n int, err error)

func DownloadFileWithRateLimitAndProgressBar(url string, wg *sync.WaitGroup) error {
	if wg != nil {
		defer wg.Done()
	}
	if *SilentMode {
		outputFunc = WriteToWgetLog
	} else {
		outputFunc = fmt.Print
	}

	var limiter *RateLimiter
	DisplayDate(true)

	if *RateLimit != "" {
		downloadSpeed, _ := ParseRateLimit(*RateLimit)
		limiter = NewLimiter(downloadSpeed)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		outputFunc("error")
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	Resp, err := http.DefaultClient.Do(req)
	if err != nil {
		outputFunc("error\n")
		return err
	}
	defer Resp.Body.Close()

	outputFunc("HTTP request sent, awaiting response... " + Resp.Status + "\n")
	if Resp.StatusCode != http.StatusOK {
		outputFunc("error\n")
		return fmt.Errorf("received non-200 status code: %s", Resp.Status)
	}

	totalSize := Resp.ContentLength
	if totalSize < 0 {
		outputFunc("Length: unspecified [text/html]\n")
	} else {
		outputFunc("Content-Length: " + strconv.FormatInt(totalSize, 10) + " (" + FormatSize(totalSize) + ")" + "\n")
	}

	outputFileName := DetermineOutputFileName(Resp, url)
	filePath := DetermineFilePath(outputFileName)
	outputFunc("saving file to:" + filePath + "\n")

	File, Any_error = os.Create(filePath)

	if Any_error != nil {
		return err
	}
	defer File.Close()

	if *SilentMode {
		outputFunc("Downloaded: " + url + "\n")
	}

	if !*SilentMode {
		if strings.Contains(Resp.Header.Get("Content-Type"), "text/html") {
			_, Any_error = io.Copy(File, Resp.Body)
			
		} else {
			bar := CreateProgressBar(totalSize)
			defer bar.Clear()
			Any_error = DownloadWithStandardProgressBar(Resp.Body,File,limiter,totalSize,bar)
		}
	} else {
		// err = writeToOutputFile(outputFileName, resp)
		_, Any_error = io.Copy(File, Resp.Body)
	}

	DisplayDate(false)
	return err
}
