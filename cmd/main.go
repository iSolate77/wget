package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
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

		baseURL, err := url.Parse(urlw)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
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
		robotsURL := baseURL.ResolveReference(&url.URL{Path: "/robots.txt"}).String()
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
		bblocks.Crawl(urlw, baseURL, discovered, client, robots)

		// Create base directory
		hostDir := path.Join(".", baseURL.Host)
		err = os.Mkdir(hostDir, 0755)
		if err != nil {
			fmt.Println("Error creating base directory:", err)
			return
		}

		// Download files
		discoveredURLs := make([]string, 0, len(discovered))
		for url := range discovered {
			discoveredURLs = append(discoveredURLs, url)
		}
		discoveredURLs = append(discoveredURLs, urlw)
		for _, url := range discoveredURLs {
			bblocks.DownloadFile(url, client, hostDir)
		}

		// Download main page if not already discovered
		mainPage := baseURL.String()
		if _, ok := discovered[mainPage]; !ok {
			bblocks.DownloadFile(mainPage, client, hostDir)
		}
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
		} else {
			urlPath := flag.Arg(0)
			if urlPath == "" {
				fmt.Println("Please provide a URL or file path")
				return
			}
			bblocks.DownloadFileWithRateLimitAndProgressBar(urlPath, nil)
		}
	}
}
