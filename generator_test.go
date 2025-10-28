package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/txtar"
)

func TestGenerateSiteE2E(t *testing.T) {
	matches, err := filepath.Glob("testdata/*.txtar")
	if err != nil {
		t.Fatal(err)
	}

	if len(matches) == 0 {
		t.Fatal("no .txtar files found in testdata/")
	}

	for _, filename := range matches {
		filename := filename
		t.Run(filepath.Base(filename), func(t *testing.T) {
			ar, err := txtar.ParseFile(filename)
			if err != nil {
				t.Fatal(err)
			}

			tmpDir := t.TempDir()
			sourceDir := filepath.Join(tmpDir, "source")
			outDir := filepath.Join(tmpDir, "out")

			if err := os.MkdirAll(sourceDir, 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(outDir, 0755); err != nil {
				t.Fatal(err)
			}

			sourceFiles := extractFiles(ar, "source/")
			wantFiles := extractFiles(ar, "want/")

			for name, content := range sourceFiles {
				path := filepath.Join(sourceDir, name)
				if err := os.WriteFile(path, content, 0644); err != nil {
					t.Fatalf("failed to write source file %s: %v", name, err)
				}
			}

			if err := generateSite(sourceDir, outDir, ""); err != nil {
				t.Fatalf("generateSite() error = %v", err)
			}

			for wantFile, wantContent := range wantFiles {
				gotPath := filepath.Join(outDir, wantFile)
				gotBytes, err := os.ReadFile(gotPath)
				if err != nil {
					t.Errorf("failed to read generated file %s: %v", wantFile, err)
					continue
				}

				gotStr := string(gotBytes)
				wantStr := strings.TrimSpace(string(wantContent))

				for _, wantSubstr := range strings.Split(wantStr, "\n") {
					wantSubstr = strings.TrimSpace(wantSubstr)
					if wantSubstr == "" {
						continue
					}
					if !strings.Contains(gotStr, wantSubstr) {
						t.Errorf("file %s missing expected content:\n  want substring: %q\n  got: %q", wantFile, wantSubstr, gotStr)
					}
				}
			}

			hiddenFile := filepath.Join(outDir, ".hidden.html")
			if _, err := os.Stat(hiddenFile); err == nil {
				t.Errorf("dotfile .hidden.md was generated but should have been ignored")
			}
		})
	}
}

func extractFiles(ar *txtar.Archive, prefix string) map[string][]byte {
	files := make(map[string][]byte)
	for _, f := range ar.Files {
		if strings.HasPrefix(f.Name, prefix) {
			name := strings.TrimPrefix(f.Name, prefix)
			files[name] = f.Data
		}
	}
	return files
}
