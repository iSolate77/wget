package bblocks

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
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

	Resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer Resp.Body.Close()

	if Resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", Resp.Status)
		return
	}

	if !robots.TestAgent(urlw, userAgent) {
		fmt.Println("Robot not allowed to crawl:", urlw)
		return
	}

	discovered[urlw] = true
	fmt.Println(urlw)

	tokenizer := html.NewTokenizer(Resp.Body)

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			switch token.Data {
			case "a", "link", "script", "img":
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
			case "style":
				for {
					tokenType := tokenizer.Next()
					if tokenType == html.ErrorToken || tokenType == html.EndTagToken && tokenizer.Token().Data == "style" {
						break
						} else if tokenType == html.TextToken {
							cssContent := tokenizer.Token().Data
							cssURLs := ExtractURLsFromCSS(cssContent, baseURL)
							for _, cssURL := range cssURLs {
								Crawl(cssURL, baseURL, discovered, client, robots)
							}
						}
					}
				}
			}
		}
	}

	func ExtractURLsFromCSS(cssContent string, baseURL *url.URL) []string {
		var urls []string

	// Regular expression to match URLs within url() declarations
	re := regexp.MustCompile(`url\(['"]?([^'"]*?)['"]?\)`)

	// Find all matches in the CSS content
	matches := re.FindAllStringSubmatch(cssContent, -1)
	fmt.Println(matches)
	for _, match := range matches {
		url := match[1] // The URL is captured in the second group
		if strings.HasPrefix(url, "(") {
			// Absolute URL
			urls = append(urls, url)
			fmt.Println(urls)
		} else {
			// Relative URL, resolve it relative to the base URL
			linkURL, err := baseURL.Parse(url)
			if err != nil {
				fmt.Println("Error resolving URL:", err)
				continue
			}
			urls = append(urls, linkURL.String())
		}
	}
	return urls
}
