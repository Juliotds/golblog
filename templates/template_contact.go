package templates

import "html/template"

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
<p style="color:#3f7398">Contact page coming soon.</p>
{{end}}`


var ContactTmpl = template.Must(template.New("contact").Parse(contactContent))
