package bblocks

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

// A function to download a file
func DownloadFile(URL_PATH string, File_name string) {

	req, err := http.NewRequest("GET", URL_PATH, nil)
	if err != nil {
		fmt.Println("error")
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error")
		return
	}
	defer resp.Body.Close()

	// Get fileName
	default_fileName := Get_filename(URL_PATH)
	var output_fileName string

	if File_name != "" {
		output_fileName = File_name
	} else if default_fileName != "" {
		output_fileName = default_fileName
	} else {
		output_fileName = "index.html"
	}

	fmt.Printf("Output will be written to %s\n", output_fileName)
	// Create file
	file, err := os.Create(output_fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	// Write to file
	// writer := bufio.NewWriter(file)
	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	if err != nil {
		fmt.Println("Error copying content to file:", err)
		return
	}

	fmt.Println("Contents saved to", output_fileName)
}
