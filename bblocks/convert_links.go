package bblocks

import (
	"io"
	"net/url"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

func ConvertHTMLLinks(input io.Reader, output io.Writer, baseURL *url.URL) error {
	doc, err := html.Parse(input)
	if err != nil {
		return errors.Wrap(err, "error parsing HTML")
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "link" || n.Data == "img") {
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
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	err = html.Render(output, doc)
	if err != nil {
		return errors.Wrap(err, "error rendering HTML")
	}

	return nil
}
