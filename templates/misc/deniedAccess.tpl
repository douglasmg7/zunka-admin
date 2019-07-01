{{template "base" .}}

{{ define "embedded-css"}} {{end}}

{{define "title"}}Acesso negado.{{end}}

{{define "content"}}
<div class="content" action="/user/change/email" method="post">
    <h2 class="title">Acesso negado</h2>
</div>
{{end}}