package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
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
  <link rel="alternate" type="application/rss+xml" title="JulioTds RSS Feed" href="/rss.xml">
  <style>
    *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

    body {
      background: #09080d;
      color: #ecdcc8;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
      line-height: 1.7;
    }

    header {
      background: #120f18;
      border-bottom: 1px solid #2e2438;
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
      color: #e85218;
      text-decoration: none;
      letter-spacing: 0.02em;
    }

    nav { display: flex; align-items: center; gap: 0.25rem; }

    nav a {
      color: #9a8070;
      text-decoration: none;
      padding: 0.4rem 0.85rem;
      border-radius: 6px;
      font-size: 0.9rem;
      transition: background 0.15s, color 0.15s;
    }

    nav a:hover { background: #2e2438; color: #ecdcc8; }

    .rss-btn {
      display: flex;
      align-items: center;
      gap: 0.35rem;
      color: #fff;
      background: #c41808;
      text-decoration: none;
      font-size: 0.85rem;
      font-weight: 600;
      padding: 0.35rem 0.75rem;
      border-radius: 6px;
      transition: background 0.15s;
      margin-left: 0.75rem;
    }

    .rss-btn:hover { background: #991005; text-decoration: none; }
    .rss-btn svg { flex-shrink: 0; }

    main {
      max-width: 760px;
      margin: 3rem auto;
      padding: 0 1.5rem;
    }

    h1, h2, h3, h4 {
      color: #f0e0cc;
      font-weight: 600;
      margin: 2rem 0 0.75rem;
      line-height: 1.3;
    }
    h1 { font-size: 2rem; }
    h2 { font-size: 1.4rem; }
    h3 { font-size: 1.15rem; }

    p { margin-bottom: 1rem; color: #c4a888; }

    a { color: #e85218; text-decoration: none; }
    a:hover { text-decoration: underline; }

    code {
      background: #1e1828;
      color: #e89828;
      padding: 0.15em 0.4em;
      border-radius: 4px;
      font-size: 0.88em;
      font-family: "JetBrains Mono", "Fira Code", monospace;
    }

    pre {
      background: #1e1828;
      border: 1px solid #2e2438;
      border-radius: 8px;
      padding: 1.25rem;
      overflow-x: auto;
      margin-bottom: 1.25rem;
    }
    pre code { background: none; padding: 0; }

    ul, ol { padding-left: 1.5rem; margin-bottom: 1rem; color: #c4a888; }
    li { margin-bottom: 0.25rem; }

    blockquote {
      border-left: 3px solid #e85218;
      padding: 0.5rem 1rem;
      margin: 1rem 0;
      color: #9a8070;
      background: #120f18;
      border-radius: 0 6px 6px 0;
    }

    hr { border: none; border-top: 1px solid #2e2438; margin: 2rem 0; }

    footer {
      background: #120f18;
      border-top: 1px solid #2e2438;
      text-align: center;
      padding: 1.5rem;
      font-size: 0.85rem;
      color: #604838;
      margin-top: 4rem;
    }

    footer a { color: #7a6050; }
    footer a:hover { color: #e85218; text-decoration: none; }

    /* Home */
    .hero {
      padding: 4rem 0 3rem;
      border-bottom: 1px solid #2e2438;
      margin-bottom: 3rem;
    }

    .hero h1 {
      font-size: 2.8rem;
      margin: 0 0 0.5rem;
      background: linear-gradient(90deg, #ff6b1a, #c81808);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }

    .hero p { font-size: 1.1rem; color: #7a6050; margin: 0; }

    .posts-section h2 {
      font-size: 1rem;
      text-transform: uppercase;
      letter-spacing: 0.1em;
      color: #604838;
      margin: 0 0 1.25rem;
    }

    .post-list { list-style: none; padding: 0; margin: 0; }
    .post-list li { border-bottom: 1px solid #1e1828; }

    .post-list a {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 1rem 0;
      color: #ecdcc8;
      text-decoration: none;
      transition: color 0.15s;
    }

    .post-list a:hover { color: #e85218; }
    .post-list .post-title { font-size: 1rem; }

    .post-list .post-slug {
      font-size: 0.8rem;
      color: #604838;
      font-family: "JetBrains Mono", monospace;
    }

    /* Blog listing */
    .blog-header {
      padding: 2.5rem 0 2rem;
      border-bottom: 1px solid #2e2438;
      margin-bottom: 2rem;
    }

    .blog-header h1 { margin: 0 0 0.25rem; font-size: 2rem; }
    .blog-header p { margin: 0; color: #7a6050; }

    .blog-list { list-style: none; padding: 0; margin: 0; }

    .blog-list li { border-bottom: 1px solid #1e1828; }

    .blog-list a {
      display: grid;
      grid-template-columns: 1fr auto;
      grid-template-rows: auto auto;
      gap: 0.25rem 1rem;
      padding: 1.25rem 0;
      color: inherit;
      text-decoration: none;
      transition: color 0.15s;
    }

    .blog-list a:hover .blog-post-title { color: #e85218; }

    .blog-post-title {
      font-size: 1rem;
      font-weight: 500;
      color: #ecdcc8;
      transition: color 0.15s;
    }

    .blog-post-date {
      font-size: 0.82rem;
      color: #604838;
      font-family: "JetBrains Mono", monospace;
      text-align: right;
      align-self: start;
    }

    .blog-post-meta {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      flex-wrap: wrap;
    }

    .blog-post-words {
      font-size: 0.78rem;
      color: #604838;
    }

    .tag {
      font-size: 0.72rem;
      padding: 0.15em 0.55em;
      border-radius: 999px;
      background: #1e1828;
      color: #e89828;
      border: 1px solid #2e2438;
    }

    .blog-list .tag { cursor: pointer; transition: background 0.15s, border-color 0.15s; }
    .blog-list .tag:hover { background: #2e2438; border-color: #e89828; }
    .blog-list .tag.tag-active { background: #e89828; color: #09080d; border-color: #e89828; }

    .tag-filter-bar {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      margin-bottom: 1.25rem;
      min-height: 1.5rem;
    }

    .tag-filter-label { font-size: 0.8rem; color: #604838; }

    .tag-clear {
      font-size: 0.75rem;
      color: #604838;
      background: none;
      border: 1px solid #2e2438;
      border-radius: 999px;
      padding: 0.1em 0.6em;
      cursor: pointer;
      transition: color 0.15s, border-color 0.15s;
    }

    .tag-clear:hover { color: #ecdcc8; border-color: #604838; }

    /* Projects */
    .projects-header {
      padding: 2.5rem 0 2rem;
      border-bottom: 1px solid #2e2438;
      margin-bottom: 2rem;
    }

    .projects-header h1 { margin: 0 0 0.25rem; font-size: 2rem; }
    .projects-header p { margin: 0; color: #7a6050; }

    .project-grid {
      display: grid;
      gap: 1rem;
    }

    .project-card {
      background: #120f18;
      border: 1px solid #2e2438;
      border-radius: 10px;
      padding: 1.25rem 1.5rem;
      text-decoration: none;
      display: block;
      transition: border-color 0.15s, background 0.15s;
    }

    .project-card:hover {
      border-color: #e85218;
      background: #1a1218;
      text-decoration: none;
    }

    .project-card-top {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      gap: 1rem;
      margin-bottom: 0.5rem;
    }

    .project-name {
      font-size: 1rem;
      font-weight: 600;
      color: #ecdcc8;
    }

    .project-year {
      font-size: 0.78rem;
      color: #604838;
      font-family: "JetBrains Mono", monospace;
      white-space: nowrap;
    }

    .project-desc {
      font-size: 0.9rem;
      color: #9a8070;
      margin: 0 0 0.85rem;
      line-height: 1.6;
    }

    .project-tags { display: flex; flex-wrap: wrap; gap: 0.4rem; }

    /* About */
    .about-header {
      padding: 3rem 0 2.5rem;
      border-bottom: 1px solid #2e2438;
      margin-bottom: 2.5rem;
    }

    .about-name {
      font-size: 2rem;
      font-weight: 700;
      color: #f0e0cc;
      margin: 0 0 0.25rem;
    }

    .about-headline {
      font-size: 1.05rem;
      color: #7a6050;
      margin: 0 0 0.75rem;
    }

    .about-location {
      font-size: 0.85rem;
      color: #604838;
      display: flex;
      align-items: center;
      gap: 0.35rem;
    }

    .about-bio {
      font-size: 1rem;
      color: #c4a888;
      line-height: 1.8;
      margin-bottom: 2rem;
    }

    .about-section-label {
      font-size: 0.8rem;
      text-transform: uppercase;
      letter-spacing: 0.1em;
      color: #604838;
      margin: 0 0 0.85rem;
    }

    .about-skills {
      display: flex;
      flex-wrap: wrap;
      gap: 0.5rem;
      margin-bottom: 2rem;
    }

    .skill-tag {
      font-size: 0.82rem;
      padding: 0.3em 0.75em;
      border-radius: 6px;
      background: #1e1828;
      color: #ecdcc8;
      border: 1px solid #2e2438;
    }

    .about-social { display: flex; flex-wrap: wrap; gap: 0.5rem; }

    .social-link {
      display: inline-flex;
      align-items: center;
      gap: 0.35rem;
      font-size: 0.88rem;
      color: #e85218;
      border: 1px solid #2e2438;
      border-radius: 6px;
      padding: 0.35rem 0.85rem;
      text-decoration: none;
      transition: background 0.15s, border-color 0.15s;
    }

    .social-link:hover {
      background: #1e1828;
      border-color: #e85218;
      text-decoration: none;
    }

    /* Comments */
    /* Contact */
    .contact-header {
      padding: 2.5rem 0 2rem;
      border-bottom: 1px solid #2e2438;
      margin-bottom: 2rem;
    }

    .contact-header h1 { margin: 0 0 0.25rem; font-size: 2rem; }
    .contact-header p { margin: 0; color: #7a6050; }

    .contact-email {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 1rem;
      color: #e85218;
      margin-bottom: 2rem;
    }

    .contact-email:hover { text-decoration: underline; }

    /* Comments */
    .comments {
      margin-top: 4rem;
      border-top: 1px solid #2e2438;
      padding-top: 2rem;
    }

    .comments-heading {
      font-size: 1rem;
      text-transform: uppercase;
      letter-spacing: 0.1em;
      color: #604838;
      margin: 0 0 1.5rem;
    }

    .comment {
      background: #120f18;
      border: 1px solid #2e2438;
      border-radius: 8px;
      padding: 1rem 1.25rem;
      margin-bottom: 1rem;
    }

    .comment-meta { display: flex; gap: 1rem; margin-bottom: 0.5rem; }
    .comment-author { font-weight: 600; color: #e85218; font-size: 0.9rem; }
    .comment-date { color: #604838; font-size: 0.85rem; }
    .comment-body { color: #c4a888; font-size: 0.95rem; margin: 0; }
    .no-comments { color: #604838; font-size: 0.9rem; margin-bottom: 1.5rem; }

    /* Comment form */
    .comment-form {
      margin-top: 2rem;
      background: #120f18;
      border: 1px solid #2e2438;
      border-radius: 8px;
      padding: 1.5rem;
    }

    .comment-form h3 {
      font-size: 1rem;
      color: #9a8070;
      margin: 0 0 1.25rem;
      font-weight: 500;
    }

    .form-group { margin-bottom: 1rem; }

    .form-optional {
      text-transform: none;
      letter-spacing: 0;
      color: #3a2818;
      font-style: italic;
    }

    .form-group label {
      display: block;
      font-size: 0.8rem;
      color: #7a6050;
      margin-bottom: 0.35rem;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    .form-group input,
    .form-group textarea {
      width: 100%;
      background: #09080d;
      border: 1px solid #2e2438;
      border-radius: 6px;
      padding: 0.6rem 0.85rem;
      color: #ecdcc8;
      font-size: 0.95rem;
      font-family: inherit;
      outline: none;
      transition: border-color 0.15s;
    }

    .form-group input:focus,
    .form-group textarea:focus { border-color: #e85218; }

    .form-group textarea { resize: vertical; min-height: 100px; }

    .form-submit {
      background: #e85218;
      color: #09080d;
      border: none;
      border-radius: 6px;
      padding: 0.6rem 1.5rem;
      font-size: 0.9rem;
      font-weight: 600;
      cursor: pointer;
      transition: background 0.15s;
    }

    .form-submit:hover { background: #ff7832; }
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
      <a class="rss-btn" href="/rss.xml">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
          <circle cx="6.18" cy="17.82" r="2.18"/>
          <path d="M4 4.44v2.83c7.03 0 12.73 5.7 12.73 12.73h2.83c0-8.59-6.97-15.56-15.56-15.56zm0 5.66v2.83c3.9 0 7.07 3.17 7.07 7.07h2.83c0-5.47-4.43-9.9-9.9-9.9z"/>
        </svg>
        RSS
      </a>
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
        <h3>Leave a comment <a href="/info" style="font-size:0.75rem;font-weight:400;color:#604838;margin-left:0.4rem;">How does this work?</a></h3>
        <form method="POST" action="/comment">
          <input type="hidden" name="slug" value="{{.Slug}}">
          <div class="form-group">
            <label for="author">Name</label>
            <input type="text" id="author" name="author" placeholder="Your name" required>
          </div>
          <div class="form-group">
            <label for="contact">Contact <span class="form-optional">(private, optional)</span></label>
            <input type="text" id="contact" name="contact" placeholder="Email or handle — only visible to the author">
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

const blogListContent = `<div class="blog-header">
  <h1>Blog</h1>
  <p>All posts, sorted by date.</p>
</div>
{{if .Posts}}
<div class="tag-filter-bar" id="tag-filter-bar" style="display:none">
  <span class="tag-filter-label">Filtering by:</span>
  <span class="tag tag-active" id="active-tag-label"></span>
  <button class="tag-clear" id="tag-clear">clear</button>
</div>
<ul class="blog-list" id="blog-list">
  {{range .Posts}}
  <li data-tags="{{range $i, $t := .Tags}}{{if $i}},{{end}}{{$t}}{{end}}">
    <a href="/blog/{{.Slug}}">
      <span class="blog-post-title">{{.Title}}</span>
      <span class="blog-post-date">{{if .Date}}{{.Date}}{{else}}&mdash;{{end}}</span>
      <span class="blog-post-meta">
        <span class="blog-post-words">{{.WordCount}} words</span>
        {{range .Tags}}<span class="tag">{{.}}</span>{{end}}
      </span>
    </a>
  </li>
  {{end}}
</ul>
<script>
(function () {
  var active = null;
  var bar = document.getElementById('tag-filter-bar');
  var label = document.getElementById('active-tag-label');
  var clear = document.getElementById('tag-clear');
  var items = document.querySelectorAll('#blog-list li');

  function filter(tag) {
    active = tag;
    label.textContent = tag;
    bar.style.display = 'flex';
    items.forEach(function (li) {
      var tags = (li.dataset.tags || '').split(',');
      li.style.display = tags.indexOf(tag) !== -1 ? '' : 'none';
    });
    document.querySelectorAll('#blog-list .tag').forEach(function (t) {
      t.classList.toggle('tag-active', t.textContent.trim() === tag);
    });
  }

  function reset() {
    active = null;
    bar.style.display = 'none';
    items.forEach(function (li) { li.style.display = ''; });
    document.querySelectorAll('#blog-list .tag').forEach(function (t) {
      t.classList.remove('tag-active');
    });
  }

  document.querySelectorAll('#blog-list .tag').forEach(function (t) {
    t.addEventListener('click', function (e) {
      e.preventDefault();
      e.stopPropagation();
      var tag = this.textContent.trim();
      tag === active ? reset() : filter(tag);
    });
  });

  clear.addEventListener('click', reset);
})();
</script>
{{else}}
  <p style="color:#604838">No posts yet.</p>
{{end}}`

const projectsContent = `<div class="projects-header">
  <h1>Projects</h1>
  <p>Things I&#39;ve built.</p>
</div>
{{if .Projects}}
<div class="project-grid">
  {{range .Projects}}
  <a class="project-card" {{if .URL}}href="{{.URL}}" target="_blank" rel="noopener"{{end}}>
    <div class="project-card-top">
      <span class="project-name">{{.Name}}</span>
      {{if .Year}}<span class="project-year">{{.Year}}</span>{{end}}
    </div>
    <p class="project-desc">{{.Description}}</p>
    <div class="project-tags">
      {{range .Tags}}<span class="tag">{{.}}</span>{{end}}
    </div>
  </a>
  {{end}}
</div>
{{else}}
  <p style="color:#604838">No projects yet.</p>
{{end}}`

const aboutContent = `{{if .}}
<div class="about-header">
  <h1 class="about-name">{{.Name}}</h1>
  <p class="about-headline">{{.Headline}}</p>
  {{if .Location}}<span class="about-location">&#128205; {{.Location}}</span>{{end}}
</div>
{{if .Bio}}<p class="about-bio">{{.Bio}}</p>{{end}}
{{if .Skills}}
<p class="about-section-label">Skills</p>
<div class="about-skills">
  {{range .Skills}}<span class="skill-tag">{{.}}</span>{{end}}
</div>
{{end}}
{{if .Social}}
<p class="about-section-label">Links</p>
<div class="about-social">
  {{range .Social}}
  <a class="social-link" href="{{.URL}}" target="_blank" rel="noopener">{{.Label}}</a>
  {{end}}
</div>
{{end}}
{{else}}
<p style="color:#604838">About page coming soon.</p>
{{end}}`

const contactContent = `{{if .}}
<div class="contact-header">
  <h1>Contact</h1>
  <p>{{.Message}}</p>
</div>
{{if .Email}}
<a class="contact-email" href="mailto:{{.Email}}">&#9993; {{.Email}}</a>
{{end}}
{{if .Social}}
<p class="about-section-label">Links</p>
<div class="about-social">
  {{range .Social}}
  <a class="social-link" href="{{.URL}}" target="_blank" rel="noopener">{{.Label}}</a>
  {{end}}
</div>
{{end}}
{{else}}
<p style="color:#604838">Contact page coming soon.</p>
{{end}}`

const infoContent = `<div class="blog-header">
  <h1>Site Info</h1>
  <p>How this blog works under the hood.</p>
</div>
<h2>Comment section</h2>
<p>
  The comment form at the bottom of each post is intentionally minimal.
  Submitting a comment sends a <code>POST</code> request to <code>/comment</code>
  with your name, message, and an optional private contact field.
</p>
<p>
  Comments are <strong>not published automatically</strong>. Each submission goes
  into a database and sits there until I review it manually. Once approved, I add
  it to the post&#39;s entry in <code>blog/comments.json</code> and regenerate the
  site, at which point it appears on the page.
</p>
<p>
  This means there will be a delay between submitting a comment and seeing it live.
  Spam, abuse, or off-topic submissions are simply discarded.
</p>
<h2>Privacy</h2>
<p>
  The optional <em>Contact</em> field in the comment form is private — it is never
  displayed publicly and is only used if I need to follow up with you directly.
</p>`

var pageTmpl = template.Must(template.New("page").Parse(htmlTemplate))
var homeTmpl = template.Must(template.New("home").Parse(homeContent))
var blogListTmpl = template.Must(template.New("blogList").Parse(blogListContent))
var projectsTmpl = template.Must(template.New("projects").Parse(projectsContent))
var aboutTmpl = template.Must(template.New("about").Parse(aboutContent))
var contactTmpl = template.Must(template.New("contact").Parse(contactContent))
var infoTmpl = template.Must(template.New("info").Parse(infoContent))

const (
	blogDir      = "blog"
	outDir       = "out"
	outBlogDir   = "out/blog"
	commentsFile = "blog/comments.json"
	projectsFile = "blog/projects.json"
	aboutFile    = "blog/about.json"
	contactFile  = "blog/contact.json"
	baseURL      = "https://juliotds.com"
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

	copied, err := copyImages(blogDir, outBlogDir)
	if err != nil {
		return fmt.Errorf("copying images: %w", err)
	}
	for _, dst := range copied {
		fmt.Printf("img  -> %s\n", dst)
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
	return renderPage(dst, contactTmpl, contact)
}

func generateInfoPage(dst string) error {
	return renderPage(dst, infoTmpl, nil)
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
	return renderPage(dst, aboutTmpl, about)
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
	return renderPage(dst, projectsTmpl, struct{ Projects []Project }{Projects: projects})
}

func generateHomePage(dst string, posts []Post) error {
	return renderPage(dst, homeTmpl, struct{ Posts []Post }{Posts: posts})
}

func generateBlogPage(dst string, posts []Post) error {
	return renderPage(dst, blogListTmpl, struct{ Posts []Post }{Posts: posts})
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
	return pageTmpl.Execute(f, PageData{Content: template.HTML(body.String())})
}

func generateRSSFeed(dst string, posts []Post) error {
	feed := rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:       "JulioTds",
			Link:        baseURL,
			Description: "Developer. Builder. Writing about code, tools, and ideas.",
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

	return pageTmpl.Execute(f, PageData{
		Content:  template.HTML(buf.String()),
		Comments: comments,
		Slug:     slug,
	})
}
