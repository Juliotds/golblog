package templates

import "html/template"

const homeContent = `<div class="hero">
  <h1>JulioTds</h1>
  <p> Writing about projects.</p>
</div>
{{if .Posts}}
<section class="posts-section">
  <h2>Recent Posts</h2>
  <ul class="post-list">
    {{range .Posts}}
    <li>
      <a href="/blog/{{.Slug}}">
        <span class="post-title">{{.Title}}</span>
        <span class="post-tags">{{range .Tags}}<span class="tag">{{.}}</span>{{end}}</span>
        {{if .Date}}<span class="blog-post-date">{{.Date}}</span>{{end}}
      </a>
    </li>
    {{end}}
  </ul>
</section>
{{end}}`


var HomeTmpl = template.Must(template.New("home").Parse(homeContent))
