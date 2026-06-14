package utils

import (
	"os"
	"strings"

	"golang.org/x/net/html"
)

func StripHTML(htmlInput string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlInput))
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && (c.Data == "script" || c.Data == "style") {
				continue
			}
			f(c)
		}
	}
	f(doc)
	return strings.TrimSpace(sb.String()), nil
}

func LoadDotenv() {
	file, err := os.ReadFile(".env")
	if err != nil {
		return
	}

	contents := string(file)
	for line := range strings.SplitSeq(contents, "\n") {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		data := strings.SplitN(line, "=", 2)
		if len(data) < 2 {
			continue
		}
		key := data[0]
		values := strings.TrimSpace(data[1])
		values = strings.Trim(values, "\"")
		os.Setenv(key, values)
	}
}
