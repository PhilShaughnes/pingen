package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
)

var wikiLinkRegex = regexp.MustCompile(`\[\[([^\]]+)\]\]`)

func convertWikiLinks(content []byte) []byte {
	return wikiLinkRegex.ReplaceAllFunc(content, func(match []byte) []byte {
		linkText := string(match[2 : len(match)-2])
		slug := makeSlug(linkText)
		return []byte(fmt.Sprintf("[%s](%s.html)", linkText, slug))
	})
}

func makeSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func parseMarkdown(content []byte) ([]byte, error) {
	content = convertWikiLinks(content)

	var buf bytes.Buffer
	md := goldmark.New()
	if err := md.Convert(content, &buf); err != nil {
		return nil, fmt.Errorf("failed to convert markdown: %w", err)
	}

	return buf.Bytes(), nil
}
