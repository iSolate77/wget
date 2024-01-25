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

// A function to download a file
func DownloadFile(URL_PATH string, File_name string, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	if *SilentMode {
		outputFunc = WriteToWgetLog
	} else {
		outputFunc = fmt.Print
	}
	if *SilentMode {
		outputFunc = WriteToWgetLog
	} else {
		outputFunc = fmt.Print
	}

	DisplayDate(true)

	req, err := http.NewRequest("GET", URL_PATH, nil)
	if err != nil {
		outputFunc("error")
		return
	}

	// Print what are you DOING!!!
	outputFunc("sending request, awaiting response... \n")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		outputFunc("error" + "\n")
		return
	}

	// print status
	outputFunc("status " + resp.Status + "\n")
	if resp.StatusCode != 200 {
		outputFunc("error" + "\n")
		return
	}

	// Print contnet size
	outputFunc("Content-Size: " + strconv.Itoa(int(resp.ContentLength)) + "\n")
	defer resp.Body.Close()

	// Get fileName
	default_fileName := Get_filename(URL_PATH)
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
	} else {
		outputFunc("saving file to:" + output_fileName + "\n")
	}
	write_to_file(output_fileName, resp)
	DisplayDate(false)
}

func write_to_file(output_fileName string, resp *http.Response) {
	// Create file
	if *New_file_path != "" {
		// File, Any_error = os.Create(*New_file_path + output_fileName)
		File, Any_error = os.Create(*New_file_path)
	} else {
		File, Any_error = os.Create(output_fileName)
	}
	if Any_error != nil {
		outputFunc("Error creating file:", Any_error)
		return
	}
	defer File.Close()

	time.Sleep(time.Second)

	// Write to file
	// writer := bufio.NewWriter(file)
	if *SilentMode {
		_, Any_error = io.Copy(File, resp.Body)
	} else {
		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			"downloading",
		)
		_, Any_error = io.Copy(io.MultiWriter(File, bar), resp.Body)
	}
	if Any_error != nil {
		outputFunc("Error copying content to file:", Any_error)
		return
	}
}
