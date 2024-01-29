package bblocks

import (
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"github.com/pkg/errors"
)

func DownloadFile(urlw string, client *http.Client, baseDir string) error {
	// Extract file name from URL
	u, err := url.Parse(urlw)
	if err != nil {
		return errors.Wrap(err, "error parsing URL")
	}

	// Get the file name from the URL path
	fileName := path.Base(u.Path)

	// If the file name doesn't have an extension, try to detect from the Content-Disposition header
	if !strings.Contains(fileName, ".") {
		resp, err := client.Head(urlw)
		if err != nil {
			return errors.Wrap(err, "error getting HEAD")
		}
		contentDisposition := resp.Header.Get("Content-Disposition")
		if contentDisposition != "" {
			_, params, err := mime.ParseMediaType(contentDisposition)
			if err == nil {
				fileName = params["filename"]
			}
		}
		contentType := resp.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "text/html") {
			fileName = "index.html"
		}
	}

	// If still no file name, use a default name
	if fileName == "" {
		fileName = "downloaded_file"
	}

	// Create directories
	filePath := path.Join(baseDir, path.Dir(u.Path))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filePath, 0755); err != nil {
			return errors.Wrap(err, "error creating directory")
		}
	}

	// Create file
	outFile, err := os.Create(path.Join(filePath, fileName))
	if err != nil {
		return errors.Wrap(err, "error creating file")
	}
	defer outFile.Close()

	// Download with progress bar
	bar := CreateProgressBar(0)

	// Define function to download file content
	downloadFunc := func() error {
		resp, err := client.Get(urlw)
		if err != nil {
			return errors.Wrap(err, "error downloading file")
		}
		defer resp.Body.Close()

		// Set total size for progress bar
		bar.ChangeMax(int(resp.ContentLength))
		totalSize := resp.ContentLength

		// Download content with progress bar
		err = DownloadWithProgressBar(resp.Body, outFile, nil, totalSize, bar)
		if err != nil {
			return errors.Wrap(err, "error downloading with progress bar")
		}

		return nil
	}

	// Execute download function
	err = downloadFunc()
	if err != nil {
		bar.Finish()
		return err
	}

	bar.Finish()
	fmt.Println("Downloaded:", fileName)
	return nil
}
