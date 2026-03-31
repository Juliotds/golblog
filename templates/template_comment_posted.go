package templates

import "html/template"

const commentPostedContent = `<div class="blog-header">
  <h1>Comment submitted</h1>
  <p>Thanks for leaving a comment.</p>
</div>
<p>
  Your comment is now pending review. Once approved it will appear on the post.
  See <a href="/info">how comments work</a> for details.
</p>
<p><a href="/blog">&larr; Back to posts</a></p>`

var CommentPostedTmpl = template.Must(template.New("commentPosted").Parse(commentPostedContent))
