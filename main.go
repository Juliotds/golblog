package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
)

const (
	blogDir = "blog"
	outDir  = "out"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	entries, err := collectMarkdownFiles(blogDir)
	if err != nil {
		return fmt.Errorf("collecting markdown files: %w", err)
	}

	for _, src := range entries {
		dst := markdownToOutputPath(src, blogDir, outDir)
		if err := convertFile(src, dst); err != nil {
			return fmt.Errorf("converting %s: %w", src, err)
		}
		fmt.Printf("%s -> %s\n", src, dst)
	}

	return nil
}

func collectMarkdownFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".md") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func markdownToOutputPath(src, srcRoot, dstRoot string) string {
	rel, _ := filepath.Rel(srcRoot, src)
	noExt := strings.TrimSuffix(rel, filepath.Ext(rel))
	return filepath.Join(dstRoot, noExt+".html")
}

func convertFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := goldmark.Convert(input, &buf); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	return os.WriteFile(dst, buf.Bytes(), 0644)
}
