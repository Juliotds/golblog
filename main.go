package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>JulioTds</title>
  <style>
    *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

    body {
      background: #0f1117;
      color: #e2e8f0;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
      line-height: 1.7;
    }

    header {
      background: #161b27;
      border-bottom: 1px solid #2d3748;
      padding: 0 2rem;
      display: flex;
      align-items: center;
      justify-content: space-between;
      height: 60px;
      position: sticky;
      top: 0;
      z-index: 10;
    }

    header .logo {
      font-size: 1.2rem;
      font-weight: 700;
      color: #a78bfa;
      text-decoration: none;
      letter-spacing: 0.02em;
    }

    nav {
      display: flex;
      gap: 0.25rem;
    }

    nav a {
      color: #94a3b8;
      text-decoration: none;
      padding: 0.4rem 0.85rem;
      border-radius: 6px;
      font-size: 0.9rem;
      transition: background 0.15s, color 0.15s;
    }

    nav a:hover {
      background: #2d3748;
      color: #e2e8f0;
    }

    main {
      max-width: 760px;
      margin: 3rem auto;
      padding: 0 1.5rem;
    }

    h1, h2, h3, h4 {
      color: #f1f5f9;
      font-weight: 600;
      margin: 2rem 0 0.75rem;
      line-height: 1.3;
    }
    h1 { font-size: 2rem; }
    h2 { font-size: 1.4rem; }
    h3 { font-size: 1.15rem; }

    p { margin-bottom: 1rem; color: #cbd5e1; }

    a { color: #a78bfa; text-decoration: none; }
    a:hover { text-decoration: underline; }

    code {
      background: #1e2535;
      color: #7dd3fc;
      padding: 0.15em 0.4em;
      border-radius: 4px;
      font-size: 0.88em;
      font-family: "JetBrains Mono", "Fira Code", monospace;
    }

    pre {
      background: #1e2535;
      border: 1px solid #2d3748;
      border-radius: 8px;
      padding: 1.25rem;
      overflow-x: auto;
      margin-bottom: 1.25rem;
    }
    pre code { background: none; padding: 0; }

    ul, ol { padding-left: 1.5rem; margin-bottom: 1rem; color: #cbd5e1; }
    li { margin-bottom: 0.25rem; }

    blockquote {
      border-left: 3px solid #a78bfa;
      padding: 0.5rem 1rem;
      margin: 1rem 0;
      color: #94a3b8;
      background: #161b27;
      border-radius: 0 6px 6px 0;
    }

    hr {
      border: none;
      border-top: 1px solid #2d3748;
      margin: 2rem 0;
    }

    footer {
      background: #161b27;
      border-top: 1px solid #2d3748;
      text-align: center;
      padding: 1.5rem;
      font-size: 0.85rem;
      color: #475569;
      margin-top: 4rem;
    }

    footer a { color: #64748b; }

    footer a:hover {
      color: #a78bfa;
      text-decoration: none;
    }

    .hero {
      padding: 4rem 0 3rem;
      border-bottom: 1px solid #2d3748;
      margin-bottom: 3rem;
    }

    .hero h1 {
      font-size: 2.8rem;
      margin: 0 0 0.5rem;
      background: linear-gradient(90deg, #a78bfa, #60a5fa);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }

    .hero p {
      font-size: 1.1rem;
      color: #64748b;
      margin: 0;
    }

    .posts-section h2 {
      font-size: 1rem;
      text-transform: uppercase;
      letter-spacing: 0.1em;
      color: #475569;
      margin: 0 0 1.25rem;
    }

    .post-list {
      list-style: none;
      padding: 0;
      margin: 0;
    }

    .post-list li { border-bottom: 1px solid #1e2535; }

    .post-list a {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 1rem 0;
      color: #e2e8f0;
      text-decoration: none;
      transition: color 0.15s;
    }

    .post-list a:hover { color: #a78bfa; }
    .post-list .post-title { font-size: 1rem; }

    .post-list .post-slug {
      font-size: 0.8rem;
      color: #475569;
      font-family: "JetBrains Mono", monospace;
    }

    /* Comments */
    .comments {
      margin-top: 4rem;
      border-top: 1px solid #2d3748;
      padding-top: 2rem;
    }

    .comments-heading {
      font-size: 1rem;
      text-transform: uppercase;
      letter-spacing: 0.1em;
      color: #475569;
      margin: 0 0 1.5rem;
    }

    .comment {
      background: #161b27;
      border: 1px solid #2d3748;
      border-radius: 8px;
      padding: 1rem 1.25rem;
      margin-bottom: 1rem;
    }

    .comment-meta {
      display: flex;
      gap: 1rem;
      margin-bottom: 0.5rem;
    }

    .comment-author {
      font-weight: 600;
      color: #a78bfa;
      font-size: 0.9rem;
    }

    .comment-date {
      color: #475569;
      font-size: 0.85rem;
    }

    .comment-body {
      color: #cbd5e1;
      font-size: 0.95rem;
      margin: 0;
    }

    .no-comments {
      color: #475569;
      font-size: 0.9rem;
      margin-bottom: 1.5rem;
    }

    /* Comment form */
    .comment-form {
      margin-top: 2rem;
      background: #161b27;
      border: 1px solid #2d3748;
      border-radius: 8px;
      padding: 1.5rem;
    }

    .comment-form h3 {
      font-size: 1rem;
      color: #94a3b8;
      margin: 0 0 1.25rem;
      font-weight: 500;
    }

    .form-group {
      margin-bottom: 1rem;
    }

    .form-group label {
      display: block;
      font-size: 0.8rem;
      color: #64748b;
      margin-bottom: 0.35rem;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    .form-group input,
    .form-group textarea {
      width: 100%;
      background: #0f1117;
      border: 1px solid #2d3748;
      border-radius: 6px;
      padding: 0.6rem 0.85rem;
      color: #e2e8f0;
      font-size: 0.95rem;
      font-family: inherit;
      outline: none;
      transition: border-color 0.15s;
    }

    .form-group input:focus,
    .form-group textarea:focus {
      border-color: #a78bfa;
    }

    .form-group textarea {
      resize: vertical;
      min-height: 100px;
    }

    .form-submit {
      background: #a78bfa;
      color: #0f1117;
      border: none;
      border-radius: 6px;
      padding: 0.6rem 1.5rem;
      font-size: 0.9rem;
      font-weight: 600;
      cursor: pointer;
      transition: background 0.15s;
    }

    .form-submit:hover { background: #c4b5fd; }
  </style>
</head>
<body>
  <header>
    <a class="logo" href="/">JulioTds</a>
    <nav>
      <a href="/">Home</a>
      <a href="/about">About</a>
      <a href="/projects">Projects</a>
      <a href="/blog">Blog</a>
      <a href="/contact">Contact</a>
    </nav>
  </header>
  <main>
    {{.Content}}
    {{if .Slug}}
    <section class="comments">
      <h2 class="comments-heading">Comments</h2>
      {{if .Comments}}
        {{range .Comments}}
        <div class="comment">
          <div class="comment-meta">
            <span class="comment-author">{{.Author}}</span>
            <span class="comment-date">{{.Date}}</span>
          </div>
          <p class="comment-body">{{.Body}}</p>
        </div>
        {{end}}
      {{else}}
        <p class="no-comments">No comments yet. Be the first!</p>
      {{end}}
      <div class="comment-form">
        <h3>Leave a comment</h3>
        <form method="POST" action="/comment">
          <input type="hidden" name="slug" value="{{.Slug}}">
          <div class="form-group">
            <label for="author">Name</label>
            <input type="text" id="author" name="author" placeholder="Your name" required>
          </div>
          <div class="form-group">
            <label for="body">Comment</label>
            <textarea id="body" name="body" placeholder="Write your comment..." required></textarea>
          </div>
          <button type="submit" class="form-submit">Post comment</button>
        </form>
      </div>
    </section>
    {{end}}
  </main>
  <footer>
    &copy; 2026 <a href="/">JulioTds</a>. All rights reserved.
  </footer>
</body>
</html>`

const homeContent = `<div class="hero">
  <h1>JulioTds</h1>
  <p>Developer. Builder. Writing about code, tools, and ideas.</p>
</div>
{{if .Posts}}
<section class="posts-section">
  <h2>Recent Posts</h2>
  <ul class="post-list">
    {{range .Posts}}
    <li>
      <a href="/blog/{{.Slug}}">
        <span class="post-title">{{.Title}}</span>
        <span class="post-slug">/blog/{{.Slug}}</span>
      </a>
    </li>
    {{end}}
  </ul>
</section>
{{end}}`

var pageTmpl = template.Must(template.New("page").Parse(htmlTemplate))
var homeTmpl = template.Must(template.New("home").Parse(homeContent))

const (
	blogDir     = "blog"
	outDir      = "out"
	commentsFile = "blog/comments.json"
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

type Post struct {
	Title string
	Slug  string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	allComments, err := loadComments(commentsFile)
	if err != nil {
		return fmt.Errorf("loading comments: %w", err)
	}

	entries, err := collectMarkdownFiles(blogDir)
	if err != nil {
		return fmt.Errorf("collecting markdown files: %w", err)
	}

	var posts []Post
	for _, src := range entries {
		post := postFromPath(src)
		dst := markdownToOutputPath(src, blogDir, outDir)
		if err := convertFile(src, dst, post.Slug, allComments[post.Slug]); err != nil {
			return fmt.Errorf("converting %s: %w", src, err)
		}
		fmt.Printf("%s -> %s\n", src, dst)
		posts = append(posts, post)
	}

	dst := filepath.Join(outDir, "index.html")
	if err := generateHomePage(dst, posts); err != nil {
		return fmt.Errorf("generating home page: %w", err)
	}
	fmt.Printf("home -> %s\n", dst)

	return nil
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

func postFromPath(src string) Post {
	rel, _ := filepath.Rel(blogDir, src)
	parts := strings.Split(filepath.ToSlash(rel), "/")
	slug := strings.TrimSuffix(parts[0], ".md")
	title := strings.ReplaceAll(slug, "-", " ")
	title = strings.Title(title)
	return Post{Title: title, Slug: slug}
}

func generateHomePage(dst string, posts []Post) error {
	var body bytes.Buffer
	if err := homeTmpl.Execute(&body, struct{ Posts []Post }{Posts: posts}); err != nil {
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

	return pageTmpl.Execute(f, PageData{Content: template.HTML(body.String())})
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

	var body bytes.Buffer
	if err := goldmark.Convert(input, &body); err != nil {
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

	return pageTmpl.Execute(f, PageData{
		Content:  template.HTML(body.String()),
		Comments: comments,
		Slug:     slug,
	})
}
