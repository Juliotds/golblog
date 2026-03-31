package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golblog/templates"

	"github.com/yuin/goldmark"
)
const (
	blogDir      = "blog"
	outDir       = "out"
	outBlogDir   = "out/blog"
	commentsFile = "blog/comments.json"
	projectsFile = "blog/projects.json"
	aboutFile    = "blog/about.json"
	contactFile  = "blog/contact.json"
	baseURL      = "https://takaguij.com"
)

type Comment struct {
	Author string `json:"author"`
	Body   string `json:"body"`
	Date   string `json:"date"`
}

type PageData struct {
	Content  template.HTML
	Comments []Comment
	Slug     string
}

type Contact struct {
	Email   string       `json:"email"`
	Message string       `json:"message"`
	Social  []SocialLink `json:"social"`
}

type SocialLink struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type About struct {
	Name     string       `json:"name"`
	Headline string       `json:"headline"`
	Bio      string       `json:"bio"`
	Location string       `json:"location"`
	Skills   []string     `json:"skills"`
	Social   []SocialLink `json:"social"`
}

type Project struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Tags        []string `json:"tags"`
	Year        string   `json:"year"`
}

type Post struct {
	Title     string
	Slug      string
	Date      string
	Tags      []string
	WordCount int
}

type rssItem struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	GUID    string   `xml:"guid"`
}

type rssChannel struct {
	XMLName     xml.Name  `xml:"channel"`
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []rssItem
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.RemoveAll(outDir); err != nil {
		return fmt.Errorf("cleaning out dir: %w", err)
	}

	allComments, err := loadComments(commentsFile)
	if err != nil {
		return fmt.Errorf("loading comments: %w", err)
	}

	projects, err := loadProjects(projectsFile)
	if err != nil {
		return fmt.Errorf("loading projects: %w", err)
	}

	about, err := loadAbout(aboutFile)
	if err != nil {
		return fmt.Errorf("loading about: %w", err)
	}

	contact, err := loadContact(contactFile)
	if err != nil {
		return fmt.Errorf("loading contact: %w", err)
	}

	entries, err := collectMarkdownFiles(blogDir)
	if err != nil {
		return fmt.Errorf("collecting markdown files: %w", err)
	}

	var posts []Post
	for _, src := range entries {
		post, err := readPost(src)
		if err != nil {
			return fmt.Errorf("reading %s: %w", src, err)
		}
		dst := markdownToOutputPath(src, blogDir, outBlogDir)
		if err := convertFile(src, dst, post.Slug, allComments[post.Slug]); err != nil {
			return fmt.Errorf("converting %s: %w", src, err)
		}
		fmt.Printf("%s -> %s\n", src, dst)
		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date > posts[j].Date
	})

	if err := generateHomePage(filepath.Join(outDir, "index.html"), posts); err != nil {
		return fmt.Errorf("generating home page: %w", err)
	}
	fmt.Printf("home -> out/index.html\n")

	if err := generateBlogPage(filepath.Join(outBlogDir, "index.html"), posts); err != nil {
		return fmt.Errorf("generating blog page: %w", err)
	}
	fmt.Printf("blog -> out/blog/index.html\n")

	if err := generateContactPage(filepath.Join(outDir, "contact", "index.html"), contact); err != nil {
		return fmt.Errorf("generating contact page: %w", err)
	}
	fmt.Printf("cont -> out/contact/index.html\n")

	if err := generateAboutPage(filepath.Join(outDir, "about", "index.html"), about); err != nil {
		return fmt.Errorf("generating about page: %w", err)
	}
	fmt.Printf("about-> out/about/index.html\n")

	if err := generateProjectsPage(filepath.Join(outDir, "projects", "index.html"), projects); err != nil {
		return fmt.Errorf("generating projects page: %w", err)
	}
	fmt.Printf("proj -> out/projects/index.html\n")

	if err := generateRSSFeed(filepath.Join(outDir, "rss.xml"), posts); err != nil {
		return fmt.Errorf("generating RSS feed: %w", err)
	}
	fmt.Printf("rss  -> out/rss.xml\n")

	if err := generateInfoPage(filepath.Join(outDir, "info", "index.html")); err != nil {
		return fmt.Errorf("generating info page: %w", err)
	}
	fmt.Printf("info -> out/info/index.html\n")

	if err := renderPage(filepath.Join(outDir, "comment-posted", "index.html"), templates.CommentPostedTmpl, nil); err != nil {
		return fmt.Errorf("generating comment-posted page: %w", err)
	}
	fmt.Printf("cmnt -> out/comment-posted/index.html\n")

	if err := renderPage(filepath.Join(outDir, "error", "index.html"), templates.ErrorTmpl, nil); err != nil {
		return fmt.Errorf("generating error page: %w", err)
	}
	fmt.Printf("err  -> out/error/index.html\n")

	if err := renderPage(filepath.Join(outDir, "invalid-operation", "index.html"), templates.InvalidOperationTmpl, nil); err != nil {
		return fmt.Errorf("generating invalid-operation page: %w", err)
	}
	fmt.Printf("invl -> out/invalid-operation/index.html\n")

	copied, err := copyImages(blogDir, outBlogDir)
	if err != nil {
		return fmt.Errorf("copying images: %w", err)
	}
	for _, dst := range copied {
		fmt.Printf("img  -> %s\n", dst)
	}

	rootStatics := []string{
		"favicon.ico",
		"favicon-16x16.png",
		"favicon-32x32.png",
		"android-chrome-192x192.png",
		"android-chrome-512x512.png",
		"apple-touch-icon.png",
		"site.webmanifest",
	}
	for _, name := range rootStatics {
		src := filepath.Join(blogDir, name)
		dst := filepath.Join(outDir, name)
		data, err := os.ReadFile(src)
		if err != nil {
			continue
		}
		if err := os.WriteFile(dst, data, 0644); err != nil {
			return fmt.Errorf("copying %s: %w", name, err)
		}
		fmt.Printf("icon -> %s\n", dst)
	}

	return nil
}

// readPost reads a markdown file, parses frontmatter, and returns a Post with metadata.
func readPost(src string) (Post, error) {
	raw, err := os.ReadFile(src)
	if err != nil {
		return Post{}, err
	}

	date, tags, body := parseFrontmatter(string(raw))

	rel, _ := filepath.Rel(blogDir, src)
	parts := strings.Split(filepath.ToSlash(rel), "/")
	slug := strings.TrimSuffix(parts[0], ".md")
	title := strings.Title(strings.ReplaceAll(slug, "-", " "))

	return Post{
		Title:     title,
		Slug:      slug,
		Date:      date,
		Tags:      tags,
		WordCount: countWords(body),
	}, nil
}

// parseFrontmatter splits YAML frontmatter from the body.
// Returns date, tags, and the remaining body content.
func parseFrontmatter(raw string) (date string, tags []string, body string) {
	if !strings.HasPrefix(raw, "---") {
		return "", nil, raw
	}
	rest := raw[3:]
	end := strings.Index(rest, "---")
	if end == -1 {
		return "", nil, raw
	}

	fm := rest[:end]
	body = strings.TrimSpace(rest[end+3:])

	for _, line := range strings.Split(fm, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "date:") {
			date = strings.TrimSpace(strings.TrimPrefix(line, "date:"))
		} else if strings.HasPrefix(line, "tags:") {
			raw := strings.TrimSpace(strings.TrimPrefix(line, "tags:"))
			raw = strings.Trim(raw, "[]")
			for _, t := range strings.Split(raw, ",") {
				if t = strings.TrimSpace(t); t != "" {
					tags = append(tags, t)
				}
			}
		}
	}
	return date, tags, body
}

func countWords(s string) int {
	return len(strings.Fields(s))
}

func loadComments(path string) (map[string][]Comment, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return map[string][]Comment{}, nil
	}
	if err != nil {
		return nil, err
	}

	var comments map[string][]Comment
	if err := json.Unmarshal(data, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func loadContact(path string) (*Contact, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var contact Contact
	if err := json.Unmarshal(data, &contact); err != nil {
		return nil, err
	}
	return &contact, nil
}

func generateContactPage(dst string, contact *Contact) error {
	return renderPage(dst, templates.ContactTmpl, contact)
}

func generateInfoPage(dst string) error {
	return renderPage(dst, templates.InfoTmpl, nil)
}

func loadAbout(path string) (*About, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var about About
	if err := json.Unmarshal(data, &about); err != nil {
		return nil, err
	}
	return &about, nil
}

func generateAboutPage(dst string, about *About) error {
	return renderPage(dst, templates.AboutTmpl, about)
}

func loadProjects(path string) ([]Project, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var projects []Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func generateProjectsPage(dst string, projects []Project) error {
	return renderPage(dst, templates.ProjectsTmpl, struct{ Projects []Project }{Projects: projects})
}

func generateHomePage(dst string, posts []Post) error {
	return renderPage(dst, templates.HomeTmpl, struct{ Posts []Post }{Posts: posts})
}

func generateBlogPage(dst string, posts []Post) error {
	return renderPage(dst, templates.BlogListTmpl, struct{ Posts []Post }{Posts: posts})
}

func renderPage(dst string, tmpl *template.Template, data any) error {
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	return templates.PageTmpl.Execute(f, PageData{Content: template.HTML(body.String())})
}

func generateRSSFeed(dst string, posts []Post) error {
	feed := rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:       "JulioTds",
			Link:        baseURL,
			Description: "Writing about projects.",
		},
	}

	for _, p := range posts {
		link := baseURL + "/blog/" + p.Slug
		feed.Channel.Items = append(feed.Channel.Items, rssItem{
			Title: p.Title,
			Link:  link,
			GUID:  link,
		})
	}

	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	return os.WriteFile(dst, append([]byte(xml.Header), data...), 0644)
}

var imageExts = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true,
	".gif": true, ".webp": true, ".svg": true, ".avif": true,
}

func copyImages(srcRoot, dstRoot string) ([]string, error) {
	var copied []string
	err := filepath.WalkDir(srcRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !imageExts[strings.ToLower(filepath.Ext(path))] {
			return nil
		}
		rel, _ := filepath.Rel(srcRoot, path)
		if !strings.Contains(rel, string(filepath.Separator)) {
			return nil // root-level files go to outDir, not outBlogDir
		}
		dst := filepath.Join(dstRoot, rel)
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.WriteFile(dst, data, 0644); err != nil {
			return err
		}
		copied = append(copied, dst)
		return nil
	})
	return copied, err
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

func convertFile(src, dst, slug string, comments []Comment) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Strip frontmatter before converting
	_, _, body := parseFrontmatter(string(input))

	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(body), &buf); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	return templates.PageTmpl.Execute(f, PageData{
		Content:  template.HTML(buf.String()),
		Comments: comments,
		Slug:     slug,
	})
}
