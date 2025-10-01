package main

import (
	"testing"
)

func TestMakeSlug(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Page Name", "page-name"},
		{"Multiple Word Page", "multiple-word-page"},
		{"with-dashes", "with-dashes"},
		{"With  Extra   Spaces", "with-extra-spaces"},
		{"Special!@#$%Characters", "special-characters"},
		{"CamelCaseText", "camelcasetext"},
		{"  leading and trailing  ", "leading-and-trailing"},
		{"123 Numbers 456", "123-numbers-456"},
		{"dash-and spaces", "dash-and-spaces"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := makeSlug(tt.input)
			if got != tt.want {
				t.Errorf("makeSlug(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestConvertWikiLinks(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single wiki link",
			input: "Check out [[Page Name]]",
			want:  "Check out [Page Name](page-name.html)",
		},
		{
			name:  "multiple wiki links",
			input: "See [[First Page]] and [[Second Page]]",
			want:  "See [First Page](first-page.html) and [Second Page](second-page.html)",
		},
		{
			name:  "wiki link with dashes",
			input: "Read [[my-page-name]]",
			want:  "Read [my-page-name](my-page-name.html)",
		},
		{
			name:  "no wiki links",
			input: "Just regular text",
			want:  "Just regular text",
		},
		{
			name:  "wiki link at start",
			input: "[[Introduction]] is important",
			want:  "[Introduction](introduction.html) is important",
		},
		{
			name:  "wiki link at end",
			input: "See the [[Conclusion]]",
			want:  "See the [Conclusion](conclusion.html)",
		},
		{
			name:  "wiki link with numbers",
			input: "Check [[Project 123]]",
			want:  "Check [Project 123](project-123.html)",
		},
		{
			name:  "empty brackets should not match",
			input: "Empty [[]] brackets",
			want:  "Empty [[]] brackets",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(convertWikiLinks([]byte(tt.input)))
			if got != tt.want {
				t.Errorf("convertWikiLinks(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "wiki link in markdown",
			input:    "# Title\n\nSee [[Other Page]]",
			contains: `<a href="other-page.html">Other Page</a>`,
		},
		{
			name:     "regular markdown link",
			input:    "[Link](https://example.com)",
			contains: `<a href="https://example.com">Link</a>`,
		},
		{
			name:     "heading conversion",
			input:    "# Heading",
			contains: "<h1>Heading</h1>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMarkdown([]byte(tt.input))
			if err != nil {
				t.Fatalf("parseMarkdown() error = %v", err)
			}
			gotStr := string(got)
			if !containsString(gotStr, tt.contains) {
				t.Errorf("parseMarkdown(%q) = %q, want to contain %q", tt.input, gotStr, tt.contains)
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
