package bblocks

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

const (
	defaultBufferSize = 32 * 1024
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		outputFunc("error\n")
		return err
	}
	defer resp.Body.Close()

	outputFunc("HTTP request sent, awaiting response... " + resp.Status + "\n")
	if resp.StatusCode != http.StatusOK {
		outputFunc("error\n")
		return fmt.Errorf("received non-200 status code: %s", resp.Status)
	}

	totalSize := resp.ContentLength
	if totalSize < 0 {
		outputFunc("Length: unspecified [text/html]\n")
	} else {
		outputFunc("Content-Length: " + strconv.FormatInt(totalSize, 10) + " ("+FormatSize(totalSize)+")" + "\n")
	}

	outputFileName := determineOutputFileName(resp, url)
	filePath := determineFilePath(outputFileName)
	outputFunc("saving file to:" + filePath + "\n")

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if *SilentMode {
		outputFunc("Downloaded: "+url+"\n")
	}

	if !*SilentMode {
		bar := createProgressBar(totalSize)
		defer bar.Clear()
		err = downloadWithProgressBar(resp.Body, file, limiter, totalSize, bar)
	} else {
		err = writeToOutputFile(outputFileName, resp)
	}

	DisplayDate(false)
	return err
}

func determineOutputFileName(resp *http.Response, url string) string {
	if *Output_name_arg_flag != "" {
		return *Output_name_arg_flag
	}
	if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return "index.html"
	}
	return Get_filename(url)
}

func determineFilePath(outputFileName string) string {
	if *New_file_path != "" {
		cleanedPath := filepath.Clean(*New_file_path)
		homeDir, _ := os.UserHomeDir()
		filePath := filepath.Join(homeDir, cleanedPath[1:], outputFileName)
		*New_file_path = filePath
		return filePath
	}
	return outputFileName
}

func createProgressBar(totalSize int64) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(
		int(totalSize),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan][1/3][reset] Writing moshable file..."),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetWidth(35),
		progressbar.OptionSetDescription(FormatSize(totalSize)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	return bar
}

func downloadWithProgressBar(body io.Reader, file *os.File, limiter *RateLimiter, totalSize int64, bar *progressbar.ProgressBar) error {
	buf := make([]byte, defaultBufferSize)
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


func writeToOutputFile(outputFileName string, resp *http.Response) error {
	 Write_to_file(outputFileName, resp)
	 return nil
}
