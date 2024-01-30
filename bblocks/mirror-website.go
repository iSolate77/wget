package bblocks

import (
	"fmt"
	"io"
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
	if *Exclude != "" {
		excludes := strings.Split(*Exclude, ",")
		for _, ex := range excludes {
			if strings.Contains(string(u.Path), ex) {
				return nil
			}
		}
	}

	// Get the file name from the URL path
	fileName := path.Base(u.Path)
	if *Reject != "" {
		rejetcs := strings.Split(*Reject, ",")
		for _, rej := range rejetcs {
			if strings.Contains(fileName, rej) {
				return nil
			}
		}
	}

	// If the file name doesn't have an extension, try to detect from the Content-Disposition header
	if !strings.Contains(fileName, ".") {
		Resp, err = client.Head(urlw)
		if err != nil {
			return errors.Wrap(err, "error getting HEAD")
		}
		contentDisposition := Resp.Header.Get("Content-Disposition")
		if contentDisposition != "" {
			_, params, err := mime.ParseMediaType(contentDisposition)
			if err == nil {
				fileName = params["filename"]
			}
		}
		contentType := Resp.Header.Get("Content-Type")
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
	OutFile, Any_error := os.Create(path.Join(filePath, fileName))
	if Any_error != nil {
		return errors.Wrap(err, "error creating file")
	}
	defer OutFile.Close()

	// Define function to download file content
	downloadFunc := func() error {
		Resp, Any_error := client.Get(urlw)
		if Any_error != nil {
			return errors.Wrap(err, "error downloading file")
		}

		defer Resp.Body.Close()

		// Set total size for progress bar
		totalSize := Resp.ContentLength
		bar := CreateProgressBar(totalSize)
		bar.ChangeMax(int(Resp.ContentLength))

		if strings.HasPrefix(Resp.Header.Get("Content-Type"), "text/html") {
			_, Any_error = io.Copy(OutFile, Resp.Body)
			if Any_error != nil {
				return Any_error
			}

		} else {
			err = DownloadWithStandardProgressBar(Resp.Body, OutFile, nil, totalSize, bar)
			if err != nil {
				return errors.Wrap(err, "error downloading with progress bar")
			}

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
