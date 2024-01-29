package bblocks

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/temoto/robotstxt"
	"golang.org/x/net/html"
)

var userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36"

func Crawl(urlw string, baseURL *url.URL, discovered map[string]bool, client *http.Client, robots *robotstxt.RobotsData) {
	if _, ok := discovered[urlw]; ok {
		return
	}

	// Parse the URL
	parsedURL, err := url.Parse(urlw)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// Check if the URL is within the same domain
	if parsedURL.Host != baseURL.Host {
		fmt.Println("Skipping external domain:", urlw)
		return
	}

	req, err := http.NewRequest("GET", urlw, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status)
		return
	}

	if !robots.TestAgent(urlw, userAgent) {
		fmt.Println("Robot not allowed to crawl:", urlw)
		return
	}

	discovered[urlw] = true
	fmt.Println(urlw)

	tokenizer := html.NewTokenizer(resp.Body)

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "a" || token.Data == "link" || token.Data == "script" || token.Data == "img" {
				for _, attr := range token.Attr {
					if (token.Data == "a" && attr.Key == "href") || (token.Data == "link" && attr.Key == "href") || (token.Data == "script" && attr.Key == "src") || (token.Data == "img" && attr.Key == "src") {
						link := attr.Val
						if strings.HasPrefix(link, "http") {
							Crawl(link, baseURL, discovered, client, robots)
						} else {
							// Resolve relative paths
							linkURL, err := baseURL.Parse(link)
							if err != nil {
								fmt.Println("Error resolving URL:", err)
								continue
							}
							Crawl(linkURL.String(), baseURL, discovered, client, robots)
						}
					}
				}
			}
		}
	}
}
