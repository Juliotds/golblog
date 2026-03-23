package templates

import "html/template"

const blogListContent = `<div class="blog-header">
  <h1>Blog</h1>
  <p>All posts, sorted by date.</p>
</div>
{{if .Posts}}
<div class="tag-filter-bar" id="tag-filter-bar" style="display:none">
  <span class="tag-filter-label">Filtering by:</span>
  <span class="tag tag-active" id="active-tag-label"></span>
  <button class="tag-clear" id="tag-clear">clear</button>
</div>
<ul class="blog-list" id="blog-list">
  {{range .Posts}}
  <li data-tags="{{range $i, $t := .Tags}}{{if $i}},{{end}}{{$t}}{{end}}">
    <a href="/blog/{{.Slug}}">
      <span class="blog-post-title">{{.Title}}</span>
      <span class="blog-post-date">{{if .Date}}{{.Date}}{{else}}&mdash;{{end}}</span>
      <span class="blog-post-meta">
        <span class="blog-post-words">{{.WordCount}} words</span>
        {{range .Tags}}<span class="tag">{{.}}</span>{{end}}
      </span>
    </a>
  </li>
  {{end}}
</ul>
<script>
(function () {
  var active = null;
  var bar = document.getElementById('tag-filter-bar');
  var label = document.getElementById('active-tag-label');
  var clear = document.getElementById('tag-clear');
  var items = document.querySelectorAll('#blog-list li');

  function filter(tag) {
    active = tag;
    label.textContent = tag;
    bar.style.display = 'flex';
    items.forEach(function (li) {
      var tags = (li.dataset.tags || '').split(',');
      li.style.display = tags.indexOf(tag) !== -1 ? '' : 'none';
    });
    document.querySelectorAll('#blog-list .tag').forEach(function (t) {
      t.classList.toggle('tag-active', t.textContent.trim() === tag);
    });
  }

  function reset() {
    active = null;
    bar.style.display = 'none';
    items.forEach(function (li) { li.style.display = ''; });
    document.querySelectorAll('#blog-list .tag').forEach(function (t) {
      t.classList.remove('tag-active');
    });
  }

  document.querySelectorAll('#blog-list .tag').forEach(function (t) {
    t.addEventListener('click', function (e) {
      e.preventDefault();
      e.stopPropagation();
      var tag = this.textContent.trim();
      tag === active ? reset() : filter(tag);
    });
  });

  clear.addEventListener('click', reset);
})();
</script>
{{else}}
  <p style="color:#3f7398">No posts yet.</p>
{{end}}`


var BlogListTmpl = template.Must(template.New("blogList").Parse(blogListContent))
