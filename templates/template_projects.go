package templates

import "html/template"

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
  <p style="color:#3f7398">No projects yet.</p>
{{end}}`


var ProjectsTmpl = template.Must(template.New("projects").Parse(projectsContent))
