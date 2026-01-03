package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Slide struct {
	Year  string `json:"year"`
	Month string `json:"month"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

func main() {
	var slides []Slide
	pattern := regexp.MustCompile(`^(\d{4})[/\\](\d{2})[/\\]([^/\\]+)[/\\]index\.html$`)
	titleRe := regexp.MustCompile(`<title>(.+?)</title>`)

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// Normalize path separators
		normalPath := filepath.ToSlash(path)
		matches := pattern.FindStringSubmatch(normalPath)
		if matches == nil {
			return nil
		}

		year, month, slug := matches[1], matches[2], matches[3]
		title := slug

		if content, err := os.ReadFile(path); err == nil {
			if m := titleRe.FindSubmatch(content); m != nil {
				title = strings.TrimSpace(string(m[1]))
			}
		}

		slides = append(slides, Slide{year, month, slug, title})
		return nil
	})

	// Sort by date descending
	sort.Slice(slides, func(i, j int) bool {
		if slides[i].Year != slides[j].Year {
			return slides[i].Year > slides[j].Year
		}
		return slides[i].Month > slides[j].Month
	})

	output, _ := json.MarshalIndent(slides, "", "  ")
	os.WriteFile("slides.json", output, 0644)

	fmt.Println("Generated slides.json:")
	fmt.Println(string(output))
}
