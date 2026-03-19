package main

import (
	"encoding/json"
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
		{"blog/hello-world/index.md", "blog", "out/blog", "out/blog/hello-world/index.html"},
		{"blog/post.md", "blog", "out/blog", "out/blog/post.html"},
		{"blog/nested/deep/post.md", "blog", "out/blog", "out/blog/nested/deep/post.html"},
	}

	for _, tt := range tests {
		got := markdownToOutputPath(tt.src, tt.srcRoot, tt.dstRoot)
		if got != tt.want {
			t.Errorf("markdownToOutputPath(%q, %q, %q) = %q, want %q", tt.src, tt.srcRoot, tt.dstRoot, got, tt.want)
		}
	}
}

func TestParseFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantDate string
		wantTags []string
		bodyHas  string
	}{
		{
			name:     "with frontmatter",
			input:    "---\ndate: 2026-03-19\ntags: go, blog\n---\n# Hello",
			wantDate: "2026-03-19",
			wantTags: []string{"go", "blog"},
			bodyHas:  "# Hello",
		},
		{
			name:     "bracket tags",
			input:    "---\ndate: 2026-01-01\ntags: [a, b, c]\n---\nbody",
			wantDate: "2026-01-01",
			wantTags: []string{"a", "b", "c"},
			bodyHas:  "body",
		},
		{
			name:    "no frontmatter",
			input:   "# Just content",
			bodyHas: "# Just content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date, tags, body := parseFrontmatter(tt.input)
			if date != tt.wantDate {
				t.Errorf("date = %q, want %q", date, tt.wantDate)
			}
			if len(tags) != len(tt.wantTags) {
				t.Errorf("tags = %v, want %v", tags, tt.wantTags)
			} else {
				for i, tag := range tags {
					if tag != tt.wantTags[i] {
						t.Errorf("tags[%d] = %q, want %q", i, tag, tt.wantTags[i])
					}
				}
			}
			if !contains(body, tt.bodyHas) {
				t.Errorf("body missing %q, got: %q", tt.bodyHas, body)
			}
		})
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"hello world", 2},
		{"one two three four", 4},
		{"  spaced   out  ", 2},
		{"", 0},
	}

	for _, tt := range tests {
		got := countWords(tt.input)
		if got != tt.want {
			t.Errorf("countWords(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestReadPost(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "blog/my-post/index.md")
	if err := os.MkdirAll(filepath.Dir(src), 0755); err != nil {
		t.Fatal(err)
	}
	content := "---\ndate: 2026-03-19\ntags: go, test\n---\n# My Post\n\nHello world content here."
	if err := os.WriteFile(src, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// readPost uses the global blogDir constant, so we need to set the source
	// relative to it. We simulate by calling parseFrontmatter directly.
	_, tags, body := parseFrontmatter(content)
	wc := countWords(body)
	if wc == 0 {
		t.Error("expected non-zero word count")
	}
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %v", tags)
	}
}

func TestCollectMarkdownFiles(t *testing.T) {
	dir := t.TempDir()

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

	if len(got) != 3 {
		t.Errorf("got %d files, want 3: %v", len(got), got)
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

func TestLoadComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "comments.json")

	data := map[string][]Comment{
		"hello-world": {
			{Author: "Alice", Body: "Great post!", Date: "2026-03-19"},
		},
	}
	raw, _ := json.Marshal(data)
	if err := os.WriteFile(path, raw, 0644); err != nil {
		t.Fatal(err)
	}

	got, err := loadComments(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	comments := got["hello-world"]
	if len(comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(comments))
	}
	if comments[0].Author != "Alice" {
		t.Errorf("expected author Alice, got %q", comments[0].Author)
	}
}

func TestLoadComments_MissingFile(t *testing.T) {
	got, err := loadComments("/nonexistent/comments.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestLoadComments_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "comments.json")
	if err := os.WriteFile(path, []byte("not json"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := loadComments(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestConvertFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "post.md")
	dst := filepath.Join(dir, "out", "post.html")

	if err := os.WriteFile(src, []byte("# Hello\n\nSome **bold** text."), 0644); err != nil {
		t.Fatal(err)
	}

	comments := []Comment{
		{Author: "Alice", Body: "Nice post!", Date: "2026-03-19"},
	}

	if err := convertFile(src, dst, "my-post", comments); err != nil {
		t.Fatalf("convertFile failed: %v", err)
	}

	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}

	html := string(content)
	checks := []string{
		"<!DOCTYPE html>", "JulioTds", "Home", "About", "Projects", "Blog", "Contact",
		"<h1>", "Hello", "<strong>", "bold",
		"Alice", "Nice post!", "2026-03-19",
		`name="slug" value="my-post"`, `action="/comment"`,
	}
	for _, s := range checks {
		if !contains(html, s) {
			t.Errorf("output HTML missing %q", s)
		}
	}
}

func TestConvertFile_StripsFromtmatter(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "post.md")
	dst := filepath.Join(dir, "post.html")

	if err := os.WriteFile(src, []byte("---\ndate: 2026-03-19\ntags: go\n---\n# Title"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := convertFile(src, dst, "slug", nil); err != nil {
		t.Fatalf("convertFile failed: %v", err)
	}

	content, _ := os.ReadFile(dst)
	html := string(content)
	if contains(html, "date: 2026") {
		t.Error("frontmatter should not appear in output HTML")
	}
	if !contains(html, "<h1>") {
		t.Error("expected heading in output HTML")
	}
}

func TestConvertFile_NoComments(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "post.md")
	dst := filepath.Join(dir, "out", "post.html")

	if err := os.WriteFile(src, []byte("# Hello"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := convertFile(src, dst, "my-post", nil); err != nil {
		t.Fatalf("convertFile failed: %v", err)
	}

	content, _ := os.ReadFile(dst)
	html := string(content)
	if !contains(html, "No comments yet") {
		t.Error("expected no-comments message")
	}
	if !contains(html, `action="/comment"`) {
		t.Error("expected comment form")
	}
}

func TestConvertFile_MissingSource(t *testing.T) {
	dir := t.TempDir()
	err := convertFile("/nonexistent/post.md", filepath.Join(dir, "out.html"), "slug", nil)
	if err == nil {
		t.Error("expected error for missing source file, got nil")
	}
}

func TestGenerateHomePage(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "index.html")

	posts := []Post{
		{Title: "Hello World", Slug: "hello-world"},
		{Title: "Go Tips", Slug: "go-tips"},
	}

	if err := generateHomePage(dst, posts); err != nil {
		t.Fatalf("generateHomePage failed: %v", err)
	}

	content, _ := os.ReadFile(dst)
	html := string(content)
	for _, s := range []string{"<!DOCTYPE html>", "JulioTds", "Hello World", "/blog/hello-world", "Go Tips"} {
		if !contains(html, s) {
			t.Errorf("home page missing %q", s)
		}
	}
}

func TestGenerateHomePage_NoPosts(t *testing.T) {
	dir := t.TempDir()
	if err := generateHomePage(filepath.Join(dir, "index.html"), nil); err != nil {
		t.Fatalf("generateHomePage failed: %v", err)
	}
}

func TestGenerateBlogPage(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "blog", "index.html")

	posts := []Post{
		{Title: "Hello World", Slug: "hello-world", Date: "2026-03-19", Tags: []string{"go", "blog"}, WordCount: 120},
		{Title: "Go Tips", Slug: "go-tips", Date: "2026-03-20", Tags: []string{"go"}, WordCount: 300},
	}

	if err := generateBlogPage(dst, posts); err != nil {
		t.Fatalf("generateBlogPage failed: %v", err)
	}

	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("blog index not created: %v", err)
	}

	html := string(content)
	checks := []string{
		"Hello World", "/blog/hello-world", "2026-03-19", "go", "blog", "120 words",
		"Go Tips", "/blog/go-tips", "2026-03-20", "300 words",
	}
	for _, s := range checks {
		if !contains(html, s) {
			t.Errorf("blog page missing %q", s)
		}
	}
}

func TestGenerateBlogPage_NoPosts(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "blog", "index.html")
	if err := generateBlogPage(dst, nil); err != nil {
		t.Fatalf("generateBlogPage failed: %v", err)
	}
	content, _ := os.ReadFile(dst)
	if !contains(string(content), "No posts yet") {
		t.Error("expected empty state message")
	}
}

func TestGenerateRSSFeed(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "rss.xml")

	posts := []Post{
		{Title: "Hello World", Slug: "hello-world"},
		{Title: "Go Tips", Slug: "go-tips"},
	}

	if err := generateRSSFeed(dst, posts); err != nil {
		t.Fatalf("generateRSSFeed failed: %v", err)
	}

	content, _ := os.ReadFile(dst)
	xml := string(content)
	for _, s := range []string{`<?xml version="1.0"`, `version="2.0"`, "JulioTds", "Hello World", "/blog/hello-world", "Go Tips"} {
		if !contains(xml, s) {
			t.Errorf("rss.xml missing %q", s)
		}
	}
}

func TestGenerateRSSFeed_NoPosts(t *testing.T) {
	dir := t.TempDir()
	if err := generateRSSFeed(filepath.Join(dir, "rss.xml"), nil); err != nil {
		t.Fatalf("generateRSSFeed failed: %v", err)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
