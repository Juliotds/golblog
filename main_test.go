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

func TestPostFromPath(t *testing.T) {
	tests := []struct {
		src       string
		wantTitle string
		wantSlug  string
	}{
		{"blog/hello-world/index.md", "Hello World", "hello-world"},
		{"blog/my-post.md", "My Post", "my-post"},
		{"blog/go-tips/index.md", "Go Tips", "go-tips"},
	}

	for _, tt := range tests {
		got := postFromPath(tt.src)
		if got.Title != tt.wantTitle {
			t.Errorf("postFromPath(%q).Title = %q, want %q", tt.src, got.Title, tt.wantTitle)
		}
		if got.Slug != tt.wantSlug {
			t.Errorf("postFromPath(%q).Slug = %q, want %q", tt.src, got.Slug, tt.wantSlug)
		}
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

	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}

	html := string(content)
	if !contains(html, "No comments yet") {
		t.Error("expected no-comments message")
	}
	if !contains(html, `action="/comment"`) {
		t.Error("expected comment form")
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

	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("rss.xml not created: %v", err)
	}

	xml := string(content)
	checks := []string{
		`<?xml version="1.0"`,
		`version="2.0"`,
		"JulioTds",
		"Hello World",
		"/blog/hello-world",
		"Go Tips",
		"/blog/go-tips",
	}
	for _, s := range checks {
		if !contains(xml, s) {
			t.Errorf("rss.xml missing %q", s)
		}
	}
}

func TestGenerateRSSFeed_NoPosts(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "rss.xml")

	if err := generateRSSFeed(dst, nil); err != nil {
		t.Fatalf("generateRSSFeed failed: %v", err)
	}

	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("rss.xml not created: %v", err)
	}

	if !contains(string(content), "JulioTds") {
		t.Error("rss.xml missing channel title")
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

	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}

	html := string(content)
	checks := []string{"<!DOCTYPE html>", "JulioTds", "Hello World", "/blog/hello-world", "Go Tips", "/blog/go-tips"}
	for _, s := range checks {
		if !contains(html, s) {
			t.Errorf("home page missing %q", s)
		}
	}
}

func TestGenerateHomePage_NoPosts(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "index.html")

	if err := generateHomePage(dst, nil); err != nil {
		t.Fatalf("generateHomePage failed: %v", err)
	}

	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}

	if !contains(string(content), "JulioTds") {
		t.Error("home page missing site title")
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
