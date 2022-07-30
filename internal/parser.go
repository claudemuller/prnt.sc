package internal

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func ScrapePicURL(htmlData string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlData))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	img := visit("", doc)

	return img, nil
}

func visit(img string, node *html.Node) string {
	if node.Type == html.ElementNode && node.Data == "img" {
		for _, d := range node.Attr {
			if d.Key == "class" && d.Val != "screenshot-image" {
				continue
			}

			if d.Key == "src" && strings.HasPrefix(d.Val, "https://i.img") {
				return d.Val
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		img = visit(img, c)
	}

	return img
}
