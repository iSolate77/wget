package bblocks

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadResource(url string) error {
	// Fetch resource from URL
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create local directory structure
	filename := Get_filename(url)
	localPath := filepath.Join(filename, url[len("http://"):])
	err = os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
	if err != nil {
		return err
	}

	// Save resource to local file
	outFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
