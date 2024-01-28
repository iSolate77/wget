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

var outputFunc func(a ...any) (n int, e error)

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

	// outputFunc("sending request, awaiting response... \n")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		outputFunc("error" + "\n")
		return err
	}

	outputFunc("HTTP request sent, awaiting response... " + resp.Status + "\n")
	if resp.StatusCode != 200 {
		outputFunc("error" + "\n")
		return err
	}

	// Print contnet size
	if resp.ContentLength < 0 {
		outputFunc("Length: unspecified [text/html]\n")
	} else {
		outputFunc("Content-Length: " + strconv.Itoa(int(resp.ContentLength)) + " ("+FormatSize((resp.ContentLength))+")" + "\n")
	}
	defer resp.Body.Close()

	default_fileName := Get_filename(url)
	var output_fileName string
	if *Output_name_arg_flag != "" {
		output_fileName = *Output_name_arg_flag
	} else {
		if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
			output_fileName = "index.html"
		} else {
			output_fileName = default_fileName
		}
	}

	if *New_file_path != "" {
		cleanedPath := filepath.Clean(*New_file_path)
		homeDir, _ := os.UserHomeDir()
		FilePath = filepath.Join(homeDir, cleanedPath[1:], output_fileName)
		*New_file_path = FilePath
		outputFunc("saving file to:" + FilePath + "\n")
		file, err = os.Create(FilePath)
		if err != nil {
			return err
		}
		defer file.Close()
	} else {
		outputFunc("saving file to:" + output_fileName + "\n")
		file, err = os.Create(output_fileName)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	if *SilentMode {
		outputFunc("Downloaded: "+url+"\n")
	}
	// Get the total size for the progress bar
	totalSize := resp.ContentLength
	var bar *progressbar.ProgressBar

	if !*SilentMode {
		// Create a custom progress bar with enhanced formatting
		bar = progressbar.NewOptions(
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
		// Download with rate limiting
		cal_rate := 32 * 1024 // Default value
		if *RateLimit != "" {
			// Assign a different value if limiter.limit is zero
			cal_rate = int(limiter.limit-10)
		}
		buf := make([]byte, cal_rate)
		for {
			// Check if the download is complete
			if totalSize <= 0 {
				break
			}
			// Read a chunk of data
			n, err := resp.Body.Read(buf)
			if err != nil && err != io.EOF {
				return err
			}
			// Reserve tokens based on the number of bytes read
			if *RateLimit != "" {
				timeRequired := limiter.Reserve(n)

				// Wait for the required time
				time.Sleep(timeRequired)
			}
			// Write the data to the output file
			_, err = file.Write(buf[:n])
			if err != nil {
				return err
			}
			// Update progress bar
			bar.Add64(int64(n))
			totalSize -= int64(n)
		}
	} else {
		Write_to_file(output_fileName, resp)
	}
	DisplayDate(false)
	return nil
}
