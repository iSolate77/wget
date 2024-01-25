package bblocks

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
)

var outputFunc func(a ...any) (n int, e error)

// A function to download a file
func DownloadFile(URL_PATH string, File_name string) {
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

	// Choosing output_fileName
	// if File_name != "" {
	// 	output_fileName = File_name
	// } else if default_fileName != "" {
	// 	output_fileName = default_fileName
	// } else {
	// 	output_fileName = "./index.html"
	// }

	// fmt.Println(resp.Header.Get("Content-Type"))

	if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		output_fileName = "index.html"
	} else {
		output_fileName = default_fileName
	}

	outputFunc("saving file to:" + output_fileName + "\n")
	write_to_file(output_fileName, resp)
	DisplayDate(false)
}

func write_to_file(output_fileName string, resp *http.Response) {
	// Create file
	file, err := os.Create(output_fileName)
	if err != nil {
		outputFunc("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write to file
	// writer := bufio.NewWriter(file)
	if *SilentMode {
		_, err = io.Copy(file, resp.Body)
	} else {
		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			"downloading",
		)
		_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	}
	if err != nil {
		outputFunc("Error copying content to file:", err)
		return
	}
}
