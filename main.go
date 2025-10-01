package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var sourceDir string
	var outDir string
	var templatePath string

	flag.StringVar(&sourceDir, "source", "./", "Source directory containing markdown files")
	flag.StringVar(&sourceDir, "s", "./", "Source directory containing markdown files (shorthand)")
	flag.StringVar(&outDir, "out", "./public", "Output directory for generated HTML")
	flag.StringVar(&outDir, "o", "./public", "Output directory for generated HTML (shorthand)")
	flag.StringVar(&templatePath, "template", "", "Custom template file (optional, uses embedded template by default)")
	flag.StringVar(&templatePath, "t", "", "Custom template file (shorthand)")
	flag.Parse()

	if err := generate(sourceDir, outDir, templatePath); err != nil {
		log.Fatalf("Error generating site: %v", err)
	}

	fmt.Printf("Site generated successfully in %s\n", outDir)
}

func generate(sourceDir, outDir, templatePath string) error {
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist: %s", sourceDir)
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return generateSite(sourceDir, outDir, templatePath)
}
