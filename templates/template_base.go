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
    @font-face {
      font-family: "JetBrains Mono";
      src: url("/fonts/JetBrainsMono-VariableFont_wght.ttf") format("truetype");
      font-weight: 100 900;
      font-style: normal;
      font-display: swap;
    }
    @font-face {
      font-family: "JetBrains Mono";
      src: url("/fonts/JetBrainsMono-Italic-VariableFont_wght.ttf") format("truetype");
      font-weight: 100 900;
      font-style: italic;
      font-display: swap;
    }
  </style>
  <style>
    /* Palette: Imaginos (1988) — gothic near-black sky, charcoal storm clouds,
       pale lightning white. Surface accent from Fire of Unknown Origin (1981)
       midnight teal. Blood-red tags retained from Tyranny and Mutation. */
    :root {
      --color-bg:            #0a0a0c;
      --color-surface:       #1a3d4a;
      --color-surface-deep:  #050507;
      --color-card-hover:    #224e5e;
      --color-border:        #2c2c32;
      --color-text:          #e8e8e0;
      --color-text-body:     #a8a898;
      --color-text-muted:    #848478;
      --color-accent:        #4a8fa0;
      --color-accent-dim:    #2e6070;
      --color-tag-bg:        #2d0a10;
      --color-tag-border:    #7a1e28;
      --color-text-alpha-15: rgba(232, 232, 224, 0.15);
      --color-text-alpha-20: rgba(232, 232, 224, 0.2);
      --color-text-alpha-25: rgba(232, 232, 224, 0.25);
      --color-text-alpha-40: rgba(232, 232, 224, 0.4);
      --color-accent-alpha-15: rgba(74, 143, 160, 0.15);
    }

    *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

    ::selection { background: var(--color-text); color: var(--color-bg); }

    body {
      background: var(--color-bg);
      color: var(--color-text);
      font-family: "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
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
      background: var(--color-surface);
      border-bottom: 1px solid var(--color-border);
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
      font-family: "JetBrains Mono", monospace;
      color: var(--color-text);
      text-decoration: none;
      letter-spacing: 0.08em;
      text-shadow: 0 0 12px var(--color-text-alpha-40);
    }

    nav { display: flex; align-items: center; gap: 0.25rem; }

    nav a {
      color: var(--color-text);
      text-decoration: none;
      padding: 0.4rem 0.85rem;
      border-radius: 6px;
      font-size: 0.9rem;
      transition: background 0.15s, color 0.15s;
    }

    nav a:hover { background: var(--color-accent-alpha-15); color: var(--color-accent); }

    .rss-btn {
      display: flex;
      align-items: center;
      gap: 0.35rem;
      color: var(--color-text);
      background: var(--color-text-alpha-15);
      text-decoration: none;
      font-size: 0.85rem;
      font-weight: 600;
      padding: 0.35rem 0.75rem;
      border-radius: 6px;
      transition: background 0.15s;
      margin-left: 0.75rem;
    }

    .rss-btn:hover { background: var(--color-accent-alpha-15); color: var(--color-accent); text-decoration: none; }
    .rss-btn svg { flex-shrink: 0; }

    main {
      max-width: 760px;
      width: 100%;
      margin: 3rem auto;
      padding: 0 1.5rem;
      flex: 1;
    }

    h1, h2, h3, h4 {
      color: var(--color-text);
      font-family: "JetBrains Mono", monospace;
      font-weight: 600;
      letter-spacing: 0.03em;
      margin: 2rem 0 0.75rem;
      line-height: 1.3;
    }
    h1 { font-size: 2.2rem; font-weight: 700; }
    h2 { font-size: 1.7rem; font-weight: 600; }
    h3 { font-size: 1.3rem; font-weight: 600; }
    h4 { font-size: 1.05rem; font-weight: 500; letter-spacing: 0.05em; }

    p { margin-bottom: 1rem; color: var(--color-text-body); }

    a { color: var(--color-text); text-decoration: none; text-shadow: 0 0 8px var(--color-text-alpha-20); }
    a:hover { text-decoration: underline; text-shadow: 0 0 12px var(--color-text-alpha-40); }

    code {
      background: var(--color-surface-deep);
      color: var(--color-text);
      padding: 0.15em 0.4em;
      border-radius: 4px;
      font-size: 0.88em;
      font-family: "JetBrains Mono", "Fira Code", monospace;
    }

    pre {
      background: var(--color-surface-deep);
      border: 1px solid var(--color-border);
      border-radius: 8px;
      padding: 1.25rem;
      overflow-x: auto;
      margin-bottom: 1.25rem;
    }
    pre code { background: none; padding: 0; }

    ul, ol { padding-left: 1.5rem; margin-bottom: 1rem; color: var(--color-text-body); }
    li { margin-bottom: 0.25rem; }

    blockquote {
      border-left: 3px solid var(--color-text);
      padding: 0.5rem 1rem;
      margin: 1rem 0;
      color: var(--color-text);
      background: var(--color-surface);
      border-radius: 0 6px 6px 0;
    }

    hr { border: none; border-top: 1px solid var(--color-border); margin: 2rem 0; }

    footer {
      background: var(--color-surface);
      border-top: 1px solid var(--color-border);
      text-align: center;
      padding: 1.5rem;
      font-size: 0.85rem;
      color: var(--color-text);
    }

    footer a { color: var(--color-text); }
    footer a:hover { color: #fff; text-decoration: none; }

    /* Home */
    .hero {
      padding: 4rem 0 3rem;
      border-bottom: 1px solid var(--color-border);
      margin-bottom: 3rem;
    }

    .hero h1 {
      font-size: 2.8rem;
      margin: 0 0 0.5rem;
      background: linear-gradient(90deg, var(--color-text), var(--color-text));
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
      filter: drop-shadow(0 0 16px var(--color-text-alpha-40));
    }

    .hero p { font-size: 1.1rem; color: var(--color-text); margin: 0; }

    .posts-section h2 {
      font-size: 1rem;
      text-transform: uppercase;
      letter-spacing: 0.1em;
      color: var(--color-text-muted);
      margin: 0 0 1.25rem;
    }

    .post-list { list-style: none; padding: 0; margin: 0; }
    .post-list li { border-bottom: 1px solid var(--color-surface-deep); }

    .post-list a {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 1rem 0;
      color: var(--color-text);
      text-decoration: none;
      transition: color 0.15s;
    }

    .post-list a:hover { color: var(--color-text); }
    .post-list .post-title { font-size: 1rem; }
    .post-list .post-tags { display: flex; gap: 0.3rem; flex-wrap: wrap; margin-left: auto; padding-right: 0.75rem; }

    .post-list .post-slug {
      font-size: 0.8rem;
      color: var(--color-text-muted);
      font-family: "JetBrains Mono", monospace;
    }

    /* Blog listing */
    .blog-header {
      padding: 2.5rem 0 2rem;
      border-bottom: 1px solid var(--color-border);
      margin-bottom: 2rem;
    }

    .blog-header h1 { margin: 0 0 0.25rem; font-size: 2rem; }
    .blog-header p { margin: 0; color: var(--color-text); }

    .blog-list { list-style: none; padding: 0; margin: 0; }

    .blog-list li { border-bottom: 1px solid var(--color-surface-deep); }

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

    .blog-list a:hover .blog-post-title { color: var(--color-text); }

    .blog-post-title {
      font-size: 1rem;
      font-weight: 500;
      color: var(--color-text);
      transition: color 0.15s;
    }

    .blog-post-date {
      font-size: 0.82rem;
      color: var(--color-text-muted);
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
      color: var(--color-text-muted);
    }

    .tag {
      font-size: 0.72rem;
      padding: 0.15em 0.55em;
      border-radius: 999px;
      background: var(--color-tag-bg);
      color: var(--color-text);
      border: 1px solid var(--color-tag-border);
    }

    .blog-list .tag { cursor: pointer; transition: background 0.15s, border-color 0.15s; }
    .blog-list .tag:hover { background: var(--color-tag-border); border-color: var(--color-text); }
    .blog-list .tag.tag-active { background: var(--color-text); color: var(--color-tag-bg); border-color: var(--color-text); }

    .tag-filter-bar {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      margin-bottom: 1.25rem;
      min-height: 1.5rem;
    }

    .tag-filter-label { font-size: 0.8rem; color: var(--color-text-muted); }

    .tag-clear {
      font-size: 0.75rem;
      color: var(--color-text-muted);
      background: none;
      border: 1px solid var(--color-border);
      border-radius: 999px;
      padding: 0.1em 0.6em;
      cursor: pointer;
      transition: color 0.15s, border-color 0.15s;
    }

    .tag-clear:hover { color: var(--color-text); border-color: var(--color-text-muted); }

    /* Projects */
    .projects-header {
      padding: 2.5rem 0 2rem;
      border-bottom: 1px solid var(--color-border);
      margin-bottom: 2rem;
    }

    .projects-header h1 { margin: 0 0 0.25rem; font-size: 2rem; }
    .projects-header p { margin: 0; color: var(--color-text); }

    .project-grid {
      display: grid;
      gap: 1rem;
    }

    .project-card {
      background: var(--color-surface);
      border: 1px solid var(--color-border);
      border-radius: 10px;
      padding: 1.25rem 1.5rem;
      text-decoration: none;
      display: block;
      transition: border-color 0.15s, background 0.15s;
    }

    .project-card:hover {
      border-color: var(--color-text);
      background: var(--color-card-hover);
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
      color: var(--color-text);
    }

    .project-year {
      font-size: 0.78rem;
      color: var(--color-text-muted);
      font-family: "JetBrains Mono", monospace;
      white-space: nowrap;
    }

    .project-desc {
      font-size: 0.9rem;
      color: var(--color-text);
      margin: 0 0 0.85rem;
      line-height: 1.6;
    }

    .project-tags { display: flex; flex-wrap: wrap; gap: 0.4rem; }

    /* About */
    .about-header {
      padding: 3rem 0 2.5rem;
      border-bottom: 1px solid var(--color-border);
      margin-bottom: 2.5rem;
    }

    .about-name {
      font-size: 2rem;
      font-weight: 700;
      color: var(--color-text);
      margin: 0 0 0.25rem;
    }

    .about-headline {
      font-size: 1.05rem;
      color: var(--color-text);
      margin: 0 0 0.75rem;
    }

    .about-location {
      font-size: 0.85rem;
      color: var(--color-text-muted);
      display: flex;
      align-items: center;
      gap: 0.35rem;
    }

    .about-bio {
      font-size: 1rem;
      color: var(--color-text-body);
      line-height: 1.8;
      margin-bottom: 2rem;
    }

    .about-section-label {
      font-size: 0.8rem;
      text-transform: uppercase;
      letter-spacing: 0.1em;
      color: var(--color-text-muted);
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
      background: var(--color-surface-deep);
      color: var(--color-text);
      border: 1px solid var(--color-border);
    }

    .about-social { display: flex; flex-wrap: wrap; gap: 0.5rem; }

    .social-link {
      display: inline-flex;
      align-items: center;
      gap: 0.35rem;
      font-size: 0.88rem;
      color: var(--color-text);
      border: 1px solid var(--color-border);
      border-radius: 6px;
      padding: 0.35rem 0.85rem;
      text-decoration: none;
      transition: background 0.15s, border-color 0.15s;
    }

    .social-link:hover {
      background: var(--color-surface-deep);
      border-color: var(--color-text);
      text-decoration: none;
    }

    /* Comments */
    /* Contact */
    .contact-header {
      padding: 2.5rem 0 2rem;
      border-bottom: 1px solid var(--color-border);
      margin-bottom: 2rem;
    }

    .contact-header h1 { margin: 0 0 0.25rem; font-size: 2rem; }
    .contact-header p { margin: 0; color: var(--color-text); }

    .contact-email {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 1rem;
      color: var(--color-text);
      margin-bottom: 2rem;
    }

    .contact-email:hover { text-decoration: underline; }

    /* Comments */
    .comments {
      margin-top: 4rem;
      border-top: 1px solid var(--color-border);
      padding-top: 2rem;
    }

    .comments-heading {
      font-size: 1rem;
      text-transform: uppercase;
      letter-spacing: 0.1em;
      color: var(--color-text-muted);
      margin: 0 0 1.5rem;
    }

    .comment {
      background: var(--color-surface);
      border: 1px solid var(--color-border);
      border-radius: 8px;
      padding: 1rem 1.25rem;
      margin-bottom: 1rem;
    }

    .comment-meta { display: flex; gap: 1rem; margin-bottom: 0.5rem; }
    .comment-author { font-weight: 600; color: var(--color-text); font-size: 0.9rem; }
    .comment-date { color: var(--color-text-muted); font-size: 0.85rem; }
    .comment-body { color: var(--color-text-body); font-size: 0.95rem; margin: 0; }
    .no-comments { color: var(--color-text-muted); font-size: 0.9rem; margin-bottom: 1.5rem; }

    /* Comment form */
    .comment-form {
      margin-top: 2rem;
      background: var(--color-surface);
      border: 1px solid var(--color-border);
      border-radius: 8px;
      padding: 1.5rem;
    }

    .comment-form h3 {
      font-size: 1rem;
      color: var(--color-text);
      margin: 0 0 1.25rem;
      font-weight: 500;
    }

    .form-group { margin-bottom: 1rem; }

    .form-optional {
      text-transform: none;
      letter-spacing: 0;
      color: var(--color-border);
      font-style: italic;
    }

    .form-group label {
      display: block;
      font-size: 0.8rem;
      color: var(--color-text);
      margin-bottom: 0.35rem;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    .form-group input,
    .form-group textarea {
      width: 100%;
      background: var(--color-bg);
      border: 1px solid var(--color-border);
      border-radius: 6px;
      padding: 0.6rem 0.85rem;
      color: var(--color-text);
      font-size: 0.95rem;
      font-family: inherit;
      outline: none;
      transition: border-color 0.15s;
    }

    .form-group input:focus,
    .form-group textarea:focus { border-color: var(--color-accent); }

    .form-group textarea { resize: vertical; min-height: 100px; }

    .form-submit {
      background: var(--color-accent);
      color: var(--color-bg);
      border: none;
      border-radius: 6px;
      padding: 0.6rem 1.5rem;
      font-size: 0.9rem;
      font-weight: 600;
      cursor: pointer;
      transition: background 0.15s;
    }

    .form-submit:hover { background: var(--color-accent-dim); }

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
        <form id="comment-form" method="POST">
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
        <script>
        document.getElementById('comment-form').addEventListener('submit', function(e) {
          e.preventDefault();
          var data = new FormData(this);
          var params = new URLSearchParams();
          data.forEach(function(value, key) { params.append(key, value); });
          this.method = 'POST';
          this.action = '/comments?' + params.toString();
          this.submit();
        });
        </script>
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
