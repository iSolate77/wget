package bblocks

import (
	"io"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func ConvertHTMLLinks(input io.Reader, output io.Writer, baseURL *url.URL) error {
	doc, err := html.Parse(input)
	if err != nil {
		return err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "link" || n.Data == "img") {
			// Convert href or src attributes
			for i, attr := range n.Attr {
				if attr.Key == "href" || attr.Key == "src" {
					u, err := url.Parse(attr.Val)
					if err != nil {
						continue
					}
					absURL := baseURL.ResolveReference(u)
					n.Attr[i].Val = absURL.String()
				}
			}
		} else if n.Data == "style" {
			// Convert URLs inside style element
			cssContent := strings.TrimSpace(getTextContent(n))
			cssContent = convertURLsInCSS(cssContent, baseURL)
			setTextContent(n, cssContent)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	err = html.Render(output, doc)
	if err != nil {
		return err
	}

	return nil
}

func getTextContent(n *html.Node) string {
	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			result += c.Data
		}
	}
	return result
}

func setTextContent(n *html.Node, content string) {
	n.FirstChild = &html.Node{Type: html.TextNode, Data: content}
	n.LastChild = n.FirstChild
}

func convertURLsInCSS(cssContent string, baseURL *url.URL) string {
	// Regular expression to match URLs within url() declarations
	re := regexp.MustCompile(`url\(['"]?([^'"]*?)['"]?\)`)

	// Replace the URLs in the CSS content
	modifiedCSS := re.ReplaceAllStringFunc(cssContent, func(match string) string {
		// Extract the URL from the match
		urlMatch := re.FindStringSubmatch(match)
		if len(urlMatch) < 2 {
			return match // No URL found, return the original match
		}
		originalURL := urlMatch[1]

		// Resolve the URL relative to the base URL
		resolvedURL := baseURL.ResolveReference(&url.URL{Path: originalURL}).String()

		// Return the modified URL wrapped in the url() declaration
		return "url('" + resolvedURL + "')"
	})

	return modifiedCSS
}

func ConvertURLs(htmlContent []byte) string {
	// Convert HTML content to string
	htmlStr := string(htmlContent)

	// Absolute path to the directory
	// Define a regex pattern to match the URL values inside url() function
	urlPattern := `url\(['"]?(.*?)['"]?\)`

	// Compile the regex pattern
	re := regexp.MustCompile(urlPattern)

	// Replace the matched URL values with the desired value
	modifiedHTML := re.ReplaceAllStringFunc(htmlStr, func(match string) string {
		// Extract the URL value from the matched string
		url := re.FindStringSubmatch(match)[1]

		// Replace the URL value with the desired absolute path
		absolutePath := "/corndog.io" + url // Replace with your working directory path
		return "url('" + absolutePath + "')"
	})

	return modifiedHTML
}
