package templates

import "html/template"

const infoContent = `<div class="blog-header">
  <h1>Site Info</h1>
  <p>How this blog works under the hood.</p>
</div>
<h2>Comment section</h2>
<p>
  The comment form at the bottom of each post is intentionally minimal.
  Submitting a comment sends a <code>POST</code> request to <code>/comments</code>
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


var InfoTmpl = template.Must(template.New("info").Parse(infoContent))
