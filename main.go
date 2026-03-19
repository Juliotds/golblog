package main

import (
	"bytes"
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
  </main>
</body>
</html>`

var pageTmpl = template.Must(template.New("page").Parse(htmlTemplate))

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

	return pageTmpl.Execute(f, struct{ Content template.HTML }{
		Content: template.HTML(body.String()),
	})
}
