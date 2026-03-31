package templates

import "html/template"

const invalidOperationContent = `<div class="blog-header">
  <h1>Invalid operation</h1>
  <p>That request is not allowed.</p>
</div>
<p>
  The action you tried to perform is invalid or missing required fields.
  Please go back and try again.
</p>
<p><a href="/">&larr; Back to home</a></p>`

var InvalidOperationTmpl = template.Must(template.New("invalidOperation").Parse(invalidOperationContent))
