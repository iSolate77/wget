package bblocks

import (
	"regexp"
)

func Get_filename(URL_PATH string) string {
	// Get the file extension from the URL
	// ext := filepath.Ext(URL_PATH)

	// Define the regular expression to capture the file name
	re := regexp.MustCompile(`\/([^\/]+)$`)

	// Find matches
	matches := re.FindStringSubmatch(URL_PATH)
	var fileName string
	// Extract the file name
	if len(matches) >= 2 {
		fileName = matches[1]
		// outputFunc("File Name: " + fileName + "\n")
	} else {
		return ""
	}

	return fileName
}
