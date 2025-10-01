package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

type Page struct {
	Title    string
	Slug     string
	FilePath string
	Content  template.HTML
}

type PageData struct {
	Title   string
	Content template.HTML
	Pages   []Page
}

func generateSite(sourceDir, outDir string) error {
	pages, err := collectPages(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to collect pages: %w", err)
	}

	tmpl, err := template.ParseFiles("templates/page.html")
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	for _, page := range pages {
		content, err := os.ReadFile(page.FilePath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", page.FilePath, err)
		}

		html, err := parseMarkdown(content)
		if err != nil {
			return fmt.Errorf("failed to parse markdown for %s: %w", page.FilePath, err)
		}

		page.Content = template.HTML(html)

		data := PageData{
			Title:   page.Title,
			Content: page.Content,
			Pages:   pages,
		}

		outPath := filepath.Join(outDir, page.Slug+".html")
		outFile, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", outPath, err)
		}
		defer outFile.Close()

		if err := tmpl.Execute(outFile, data); err != nil {
			return fmt.Errorf("failed to execute template for %s: %w", page.Title, err)
		}

		fmt.Printf("Generated: %s -> %s\n", page.FilePath, outPath)
	}

	return nil
}

func collectPages(sourceDir string) ([]Page, error) {
	var pages []Page

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") && path != sourceDir {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		title := strings.TrimSuffix(info.Name(), ".md")
		slug := makeSlug(title)

		pages = append(pages, Page{
			Title:    title,
			Slug:     slug,
			FilePath: path,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return pages, nil
}
