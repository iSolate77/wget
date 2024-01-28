package bblocks

import (
	"io"

	"golang.org/x/net/html"
)

func ExtractLinks(body io.Reader) ([]string, error) {
	links := make([]string, 0)

	// Parse HTML
	doc, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	// Define a recursive function to traverse the HTML tree
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			// Extract href attribute from <a> tags
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
					break
				}
			}
		}
		// Recursively process child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	// Start traversing the HTML tree
	extract(doc)

	return links, nil
}
