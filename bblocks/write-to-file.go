package bblocks

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

func Write_to_files(output_fileName string, resp *http.Response) {
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
