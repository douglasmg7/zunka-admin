{{template "base" .}}

{{ define "embedded-css"}} 
<style type="text/css">
    h2 {
        text-align: center;
    }
</style>
{{end}}

{{define "title"}}Acesso negado{{end}}

{{define "content"}}
<div class="content" action="/ns/user/change/email" method="post">
    <h2 class="title">Acesso negado</h2>
</div>
{{end}}