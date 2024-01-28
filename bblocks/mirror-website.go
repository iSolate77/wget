package bblocks

import (
	"net/http"
)

func MirrorWebsite(baseUrl string, resources *[]string) error {
	// Fetch HTML content of the current URL
	resp, err := http.Get(baseUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Extract links from HTML content
	links, errs := ExtractLinks(resp.Body)
	if errs != nil {
		return errs
	}
	// Append discovered links to resources array
	for _, link := range links {
		if !contains(*resources, link) {
			*resources = append(*resources, link)
		}
	}

	// Download resources recursively

	// Download resource if it's not already downloaded
	for _, resource := range *resources {
		DownloadResource(resource)
	}

	return nil
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
