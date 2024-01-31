package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"
	"wget/bblocks"

	"github.com/temoto/robotstxt"
)

var wg sync.WaitGroup

func main() {
	flag.Parse()

	if *bblocks.MirrorMode {
		urlw := flag.Arg(0)
		if urlw == "" {
			fmt.Println("Please provide a URL to mirror")
			return
		}

		bblocks.BaseUrl, bblocks.Any_error = url.Parse(urlw)
		if bblocks.Any_error != nil {
			fmt.Println("Error parsing URL:", bblocks.Any_error)
			return
		}

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Allow automatic redirects
				return nil
			},
			Timeout: 10 * time.Second,
		}

		// Fetch robots.txt and parse
		robotsURL := bblocks.BaseUrl.ResolveReference(&url.URL{Path: "/robots.txt"}).String()
		robotsResp, err := client.Get(robotsURL)
		if err != nil {
			fmt.Println("Error fetching robots.txt:", err)
			return
		}
		defer robotsResp.Body.Close()
		robots, err := robotstxt.FromResponse(robotsResp)
		if err != nil {
			fmt.Println("Error parsing robots.txt:", err)
			return
		}

		discovered := make(map[string]bool)
		bblocks.Crawl(urlw, bblocks.BaseUrl, discovered, client, robots)

		// Create base directory
		hostDir := path.Join(".", bblocks.BaseUrl.Host)
		err = os.Mkdir(hostDir, 0755)
		if err != nil {
			fmt.Println("Error creating base directory:", err)
			return
		}

		// Remove duplicate URLs
		uniqueDiscovered := make(map[string]bool)
		for url := range discovered {
			// Check if the URL is not already in the uniqueDiscovered map
			if _, ok := uniqueDiscovered[url]; !ok {
				uniqueDiscovered[url] = true
			}
		}

		// Convert unique URLs back to a slice
		discoveredURLs := make([]string, 0, len(uniqueDiscovered))
		for url := range uniqueDiscovered {
			discoveredURLs = append(discoveredURLs, url)
		}

		// Download files
		for _, url := range discoveredURLs {
			bblocks.DownloadFile(url, client, hostDir)
		}
		if *bblocks.ConvertMode && strings.HasPrefix(bblocks.Resp.Header.Get("Content-Type"), "text/html") {
			htmlContent, _ := ioutil.ReadFile("corndog.io/index.html")
			modifiedHTML := bblocks.ConvertURLs(htmlContent)
			// ConvertHTMLLinks(Resp.Body, File, BaseUrl)
			err = ioutil.WriteFile("corndog.io/index.html", []byte(modifiedHTML), 0644)
			if err != nil {
				fmt.Println("Error writing modified HTML file:", err)
				os.Exit(1)
			}
		}

		os.Remove("wget-log.txt")
	} else {
		if *bblocks.SilentMode {
			fmt.Println("output will be written to wget-log")
		} else {
			os.Remove("wget-log.txt")
		}
		if *bblocks.AsyncFileInput != "" {
			links, err := bblocks.GetLinksFromFile()
			if err != nil {
				fmt.Println("Error getting links from file:", err)
				return
			}
			for _, link := range links {
				wg.Add(1)
				go bblocks.DownloadFileWithRateLimitAndProgressBar(link, &wg)
			}
			wg.Wait()
			// os.Remove("wget-log.txt")
		} else {
			urlPath := flag.Args()
			for _, link := range urlPath {
				bblocks.DownloadFileWithRateLimitAndProgressBar(link, nil)
			}
			if len(urlPath) == 0 {
				fmt.Println("Please provide a URL or file path")
				return
			}
		}
	}
}
