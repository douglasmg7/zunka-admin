{{template "base" .}}
{{ define "embedded-css"}} {{end}}
{{define "title"}} Alunos {{end}}

{{define "header"}}
<div class="header">
    <h1>Alunos</h1>
    <h4>Relação de todos os alunos</h4>
</div>
{{end}}

{{define "content"}}
<div class="content">
  {{range .Students}}
    <h3 class="subtitle">
      <a href="/student/id/{{.Id}}">{{.Name}}</a>
    </h3>
  {{end}}
</div>
{{end}}