package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMarkdownToOutputPath(t *testing.T) {
	tests := []struct {
		src     string
		srcRoot string
		dstRoot string
		want    string
	}{
		{"blog/hello-world/index.md", "blog", "out", "out/hello-world/index.html"},
		{"blog/post.md", "blog", "out", "out/post.html"},
		{"blog/nested/deep/post.md", "blog", "out", "out/nested/deep/post.html"},
	}

	for _, tt := range tests {
		got := markdownToOutputPath(tt.src, tt.srcRoot, tt.dstRoot)
		if got != tt.want {
			t.Errorf("markdownToOutputPath(%q, %q, %q) = %q, want %q", tt.src, tt.srcRoot, tt.dstRoot, got, tt.want)
		}
	}
}

func TestCollectMarkdownFiles(t *testing.T) {
	dir := t.TempDir()

	// Create a mix of .md and non-.md files across nested dirs
	files := map[string]string{
		"post.md":             "# Post",
		"draft.txt":           "not markdown",
		"sub/nested.md":       "# Nested",
		"sub/deep/article.md": "# Article",
		"sub/image.png":       "",
	}
	for rel, content := range files {
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	got, err := collectMarkdownFiles(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Only .md files should be returned
	wantCount := 3
	if len(got) != wantCount {
		t.Errorf("got %d files, want %d: %v", len(got), wantCount, got)
	}
	for _, f := range got {
		if filepath.Ext(f) != ".md" {
			t.Errorf("non-.md file returned: %s", f)
		}
	}
}

func TestCollectMarkdownFiles_EmptyDir(t *testing.T) {
	dir := t.TempDir()

	got, err := collectMarkdownFiles(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected no files, got %v", got)
	}
}

func TestCollectMarkdownFiles_MissingDir(t *testing.T) {
	_, err := collectMarkdownFiles("/nonexistent/path")
	if err == nil {
		t.Error("expected error for missing directory, got nil")
	}
}

func TestConvertFile(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "post.md")
	dst := filepath.Join(dir, "out", "post.html")

	if err := os.WriteFile(src, []byte("# Hello\n\nSome **bold** text."), 0644); err != nil {
		t.Fatal(err)
	}

	if err := convertFile(src, dst); err != nil {
		t.Fatalf("convertFile failed: %v", err)
	}

	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}

	html := string(content)
	if html == "" {
		t.Error("output HTML is empty")
	}

	checks := []string{"<h1>", "Hello", "<strong>", "bold"}
	for _, s := range checks {
		if !contains(html, s) {
			t.Errorf("output HTML missing %q\ngot: %s", s, html)
		}
	}
}

func TestConvertFile_MissingSource(t *testing.T) {
	dir := t.TempDir()
	err := convertFile("/nonexistent/post.md", filepath.Join(dir, "out.html"))
	if err == nil {
		t.Error("expected error for missing source file, got nil")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
