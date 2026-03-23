package templates

import "html/template"

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
<p style="color:#3f7398">About page coming soon.</p>
{{end}}`


var AboutTmpl = template.Must(template.New("about").Parse(aboutContent))
