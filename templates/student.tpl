{{template "base" .}}
{{ define "embedded-css"}} {{end}}
{{define "title"}} Aluno {{end}}

{{define "header"}}
<div class="header">
    <h1>{{.Name}}</h1>
    <h4>{{.Email}}</h4>
</div>
{{end}}

{{define "content"}}
<div class="content">
  <div class="container">
    <p class="subtitle">{{.Name}}</p>
    <p class="subtitle">{{.Email}}</p>
    <p class="subtitle">{{.Mobile}}</p>
  </div>
</div>
{{end}}