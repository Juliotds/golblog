package templates

import "html/template"

const errorContent = `<div class="blog-header">
  <h1>Something went wrong</h1>
  <p>An unexpected error occurred.</p>
</div>
<p>
  Please try again. If the problem persists feel free to reach out via the <a href="/contact">contact page</a>.
</p>
<p><a href="/">&larr; Back to home</a></p>`

var ErrorTmpl = template.Must(template.New("error").Parse(errorContent))
