package templates

import "html/template"

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>JulioTds</title>
  <link rel="icon" type="image/x-icon" href="/favicon.ico">
  <link rel="alternate" type="application/rss+xml" title="JulioTds RSS Feed" href="/rss.xml">
  <style>
    *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

    ::selection { background: #cad1ce; color: #011140; }

    body {
      background: #011140;
      color: #cad1ce;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
      line-height: 1.7;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
    }

    /* Scanline overlay */
    body::after {
      content: '';
      position: fixed;
      inset: 0;
      background: repeating-linear-gradient(
        0deg,
        transparent,
        transparent 3px,
        rgba(0, 0, 0, 0.06) 3px,
        rgba(0, 0, 0, 0.06) 4px
      );
      pointer-events: none;
      z-index: 9999;
    }

    header {
      background: #18608C;
      border-bottom: 1px solid #0a2a6e;
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
      font-family: "JetBrains Mono", "Courier New", monospace;
      color: #cad1ce;
      text-decoration: none;
      letter-spacing: 0.08em;
      text-shadow: 0 0 12px rgba(202, 209, 206, 0.4);
    }

    nav { display: flex; align-items: center; gap: 0.25rem; }

    nav a {
      color: #cad1ce;
      text-decoration: none;
      padding: 0.4rem 0.85rem;
      border-radius: 6px;
      font-size: 0.9rem;
      transition: background 0.15s, color 0.15s;
    }

    nav a:hover { background: rgba(202,209,206,0.15); color: #fff; }

    .rss-btn {
      display: flex;
      align-items: center;
      gap: 0.35rem;
      color: #cad1ce;
      background: rgba(202,209,206,0.15);
      text-decoration: none;
      font-size: 0.85rem;
      font-weight: 600;
      padding: 0.35rem 0.75rem;
      border-radius: 6px;
      transition: background 0.15s;
      margin-left: 0.75rem;
    }

    .rss-btn:hover { background: rgba(202,209,206,0.25); text-decoration: none; }
    .rss-btn svg { flex-shrink: 0; }

    main {
      max-width: 760px;
      width: 100%;
      margin: 3rem auto;
      padding: 0 1.5rem;
      flex: 1;
    }

    h1, h2, h3, h4 {
      color: #cad1ce;
      font-family: "JetBrains Mono", "Courier New", monospace;
      font-weight: 600;
      letter-spacing: 0.03em;
      margin: 2rem 0 0.75rem;
      line-height: 1.3;
    }
    h1 { font-size: 2rem; }
    h2 { font-size: 1.4rem; }
    h3 { font-size: 1.15rem; }

    p { margin-bottom: 1rem; color: #8ec5d8; }

    a { color: #cad1ce; text-decoration: none; text-shadow: 0 0 8px rgba(202, 209, 206, 0.2); }
    a:hover { text-decoration: underline; text-shadow: 0 0 12px rgba(202, 209, 206, 0.4); }

    code {
      background: #000a2e;
      color: #cad1ce;
      padding: 0.15em 0.4em;
      border-radius: 4px;
      font-size: 0.88em;
      font-family: "JetBrains Mono", "Fira Code", monospace;
    }

    pre {
      background: #000a2e;
      border: 1px solid #0a2a6e;
      border-radius: 8px;
      padding: 1.25rem;
      overflow-x: auto;
      margin-bottom: 1.25rem;
    }
    pre code { background: none; padding: 0; }

    ul, ol { padding-left: 1.5rem; margin-bottom: 1rem; color: #8ec5d8; }
    li { margin-bottom: 0.25rem; }

    blockquote {
      border-left: 3px solid #cad1ce;
      padding: 0.5rem 1rem;
      margin: 1rem 0;
      color: #cad1ce;
      background: #18608C;
      border-radius: 0 6px 6px 0;
    }

    hr { border: none; border-top: 1px solid #0a2a6e; margin: 2rem 0; }

    footer {
      background: #18608C;
      border-top: 1px solid #0a2a6e;
      text-align: center;
      padding: 1.5rem;
      font-size: 0.85rem;
      color: #cad1ce;
    }

    footer a { color: #cad1ce; }
    footer a:hover { color: #fff; text-decoration: none; }

    /* Home */
    .hero {
      padding: 4rem 0 3rem;
      border-bottom: 1px solid #0a2a6e;
      margin-bottom: 3rem;
    }

    .hero h1 {
      font-size: 2.8rem;
      margin: 0 0 0.5rem;
      background: linear-gradient(90deg, #cad1ce, #cad1ce);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
      filter: drop-shadow(0 0 16px rgba(202, 209, 206, 0.4));
    }

    .hero p { font-size: 1.1rem; color: #cad1ce; margin: 0; }

    .posts-section h2 {
      font-size: 1rem;
      text-transform: uppercase;
      letter-spacing: 0.1em;
      color: #6aaac0;
      margin: 0 0 1.25rem;
    }

    .post-list { list-style: none; padding: 0; margin: 0; }
    .post-list li { border-bottom: 1px solid #000a2e; }

    .post-list a {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 1rem 0;
      color: #cad1ce;
      text-decoration: none;
      transition: color 0.15s;
    }

    .post-list a:hover { color: #cad1ce; }
    .post-list .post-title { font-size: 1rem; }
    .post-list .post-tags { display: flex; gap: 0.3rem; flex-wrap: wrap; margin-left: auto; padding-right: 0.75rem; }

    .post-list .post-slug {
      font-size: 0.8rem;
      color: #6aaac0;
      font-family: "JetBrains Mono", monospace;
    }

    /* Blog listing */
    .blog-header {
      padding: 2.5rem 0 2rem;
      border-bottom: 1px solid #0a2a6e;
      margin-bottom: 2rem;
    }

    .blog-header h1 { margin: 0 0 0.25rem; font-size: 2rem; }
    .blog-header p { margin: 0; color: #cad1ce; }

    .blog-list { list-style: none; padding: 0; margin: 0; }

    .blog-list li { border-bottom: 1px solid #000a2e; }

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

    .blog-list a:hover .blog-post-title { color: #cad1ce; }

    .blog-post-title {
      font-size: 1rem;
      font-weight: 500;
      color: #cad1ce;
      transition: color 0.15s;
    }

    .blog-post-date {
      font-size: 0.82rem;
      color: #6aaac0;
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
      color: #6aaac0;
    }

    .tag {
      font-size: 0.72rem;
      padding: 0.15em 0.55em;
      border-radius: 999px;
      background: #000a2e;
      color: #cad1ce;
      border: 1px solid #0a2a6e;
    }

    .blog-list .tag { cursor: pointer; transition: background 0.15s, border-color 0.15s; }
    .blog-list .tag:hover { background: #0a2a6e; border-color: #cad1ce; }
    .blog-list .tag.tag-active { background: #cad1ce; color: #011140; border-color: #cad1ce; }

    .tag-filter-bar {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      margin-bottom: 1.25rem;
      min-height: 1.5rem;
    }

    .tag-filter-label { font-size: 0.8rem; color: #6aaac0; }

    .tag-clear {
      font-size: 0.75rem;
      color: #6aaac0;
      background: none;
      border: 1px solid #0a2a6e;
      border-radius: 999px;
      padding: 0.1em 0.6em;
      cursor: pointer;
      transition: color 0.15s, border-color 0.15s;
    }

    .tag-clear:hover { color: #cad1ce; border-color: #6aaac0; }

    /* Projects */
    .projects-header {
      padding: 2.5rem 0 2rem;
      border-bottom: 1px solid #0a2a6e;
      margin-bottom: 2rem;
    }

    .projects-header h1 { margin: 0 0 0.25rem; font-size: 2rem; }
    .projects-header p { margin: 0; color: #cad1ce; }

    .project-grid {
      display: grid;
      gap: 1rem;
    }

    .project-card {
      background: #18608C;
      border: 1px solid #0a2a6e;
      border-radius: 10px;
      padding: 1.25rem 1.5rem;
      text-decoration: none;
      display: block;
      transition: border-color 0.15s, background 0.15s;
    }

    .project-card:hover {
      border-color: #cad1ce;
      background: #051a5a;
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
      color: #cad1ce;
    }

    .project-year {
      font-size: 0.78rem;
      color: #6aaac0;
      font-family: "JetBrains Mono", monospace;
      white-space: nowrap;
    }

    .project-desc {
      font-size: 0.9rem;
      color: #cad1ce;
      margin: 0 0 0.85rem;
      line-height: 1.6;
    }

    .project-tags { display: flex; flex-wrap: wrap; gap: 0.4rem; }

    /* About */
    .about-header {
      padding: 3rem 0 2.5rem;
      border-bottom: 1px solid #0a2a6e;
      margin-bottom: 2.5rem;
    }

    .about-name {
      font-size: 2rem;
      font-weight: 700;
      color: #cad1ce;
      margin: 0 0 0.25rem;
    }

    .about-headline {
      font-size: 1.05rem;
      color: #cad1ce;
      margin: 0 0 0.75rem;
    }

    .about-location {
      font-size: 0.85rem;
      color: #6aaac0;
      display: flex;
      align-items: center;
      gap: 0.35rem;
    }

    .about-bio {
      font-size: 1rem;
      color: #8ec5d8;
      line-height: 1.8;
      margin-bottom: 2rem;
    }

    .about-section-label {
      font-size: 0.8rem;
      text-transform: uppercase;
      letter-spacing: 0.1em;
      color: #6aaac0;
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
      background: #000a2e;
      color: #cad1ce;
      border: 1px solid #0a2a6e;
    }

    .about-social { display: flex; flex-wrap: wrap; gap: 0.5rem; }

    .social-link {
      display: inline-flex;
      align-items: center;
      gap: 0.35rem;
      font-size: 0.88rem;
      color: #cad1ce;
      border: 1px solid #0a2a6e;
      border-radius: 6px;
      padding: 0.35rem 0.85rem;
      text-decoration: none;
      transition: background 0.15s, border-color 0.15s;
    }

    .social-link:hover {
      background: #000a2e;
      border-color: #cad1ce;
      text-decoration: none;
    }

    /* Comments */
    /* Contact */
    .contact-header {
      padding: 2.5rem 0 2rem;
      border-bottom: 1px solid #0a2a6e;
      margin-bottom: 2rem;
    }

    .contact-header h1 { margin: 0 0 0.25rem; font-size: 2rem; }
    .contact-header p { margin: 0; color: #cad1ce; }

    .contact-email {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 1rem;
      color: #cad1ce;
      margin-bottom: 2rem;
    }

    .contact-email:hover { text-decoration: underline; }

    /* Comments */
    .comments {
      margin-top: 4rem;
      border-top: 1px solid #0a2a6e;
      padding-top: 2rem;
    }

    .comments-heading {
      font-size: 1rem;
      text-transform: uppercase;
      letter-spacing: 0.1em;
      color: #6aaac0;
      margin: 0 0 1.5rem;
    }

    .comment {
      background: #18608C;
      border: 1px solid #0a2a6e;
      border-radius: 8px;
      padding: 1rem 1.25rem;
      margin-bottom: 1rem;
    }

    .comment-meta { display: flex; gap: 1rem; margin-bottom: 0.5rem; }
    .comment-author { font-weight: 600; color: #cad1ce; font-size: 0.9rem; }
    .comment-date { color: #6aaac0; font-size: 0.85rem; }
    .comment-body { color: #8ec5d8; font-size: 0.95rem; margin: 0; }
    .no-comments { color: #6aaac0; font-size: 0.9rem; margin-bottom: 1.5rem; }

    /* Comment form */
    .comment-form {
      margin-top: 2rem;
      background: #18608C;
      border: 1px solid #0a2a6e;
      border-radius: 8px;
      padding: 1.5rem;
    }

    .comment-form h3 {
      font-size: 1rem;
      color: #cad1ce;
      margin: 0 0 1.25rem;
      font-weight: 500;
    }

    .form-group { margin-bottom: 1rem; }

    .form-optional {
      text-transform: none;
      letter-spacing: 0;
      color: #0a2a6e;
      font-style: italic;
    }

    .form-group label {
      display: block;
      font-size: 0.8rem;
      color: #cad1ce;
      margin-bottom: 0.35rem;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    .form-group input,
    .form-group textarea {
      width: 100%;
      background: #011140;
      border: 1px solid #0a2a6e;
      border-radius: 6px;
      padding: 0.6rem 0.85rem;
      color: #cad1ce;
      font-size: 0.95rem;
      font-family: inherit;
      outline: none;
      transition: border-color 0.15s;
    }

    .form-group input:focus,
    .form-group textarea:focus { border-color: #cad1ce; }

    .form-group textarea { resize: vertical; min-height: 100px; }

    .form-submit {
      background: #cad1ce;
      color: #011140;
      border: none;
      border-radius: 6px;
      padding: 0.6rem 1.5rem;
      font-size: 0.9rem;
      font-weight: 600;
      cursor: pointer;
      transition: background 0.15s;
    }

    .form-submit:hover { background: #cad1ce; }

    @media (max-width: 700px) {
      header {
        height: auto;
        padding: 0.6rem 1rem;
        flex-wrap: wrap;
        gap: 0.4rem;
      }

      nav {
        flex-wrap: wrap;
        gap: 0.1rem;
      }

      nav a {
        font-size: 0.82rem;
        padding: 0.3rem 0.6rem;
      }

      .rss-btn {
        font-size: 0.78rem;
        padding: 0.28rem 0.55rem;
        margin-left: 0;
      }

      main {
        margin: 1.75rem auto;
        padding: 0 1rem;
      }

      h1 { font-size: 1.6rem; }
      h2 { font-size: 1.2rem; }
      h3 { font-size: 1rem; }

      .hero {
        padding: 2rem 0 1.5rem;
      }

      .hero h1 {
        font-size: 2rem;
      }

      .post-list .post-slug { display: none; }

      .blog-list a {
        grid-template-columns: 1fr;
      }

      .blog-post-date {
        text-align: left;
      }

      .project-card {
        padding: 1rem 1.1rem;
      }

      .about-name { font-size: 1.5rem; }

      .comment-form {
        padding: 1rem;
      }
    }

    @media (max-width: 420px) {
      .hero h1 { font-size: 1.6rem; }

      nav a { font-size: 0.78rem; padding: 0.25rem 0.5rem; }
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
        <h3>Leave a comment <a href="/info" style="font-size:0.75rem;font-weight:400;color:#3f7398;margin-left:0.4rem;">How does this work?</a></h3>
        <form method="POST" action="/comments">
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


var PageTmpl = template.Must(template.New("page").Parse(htmlTemplate))
